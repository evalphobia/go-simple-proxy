package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"
)

var (
	// under score variables are used as flag value.
	_verbose bool
	_debug   bool
	_timeout string

	timeout time.Duration
)

var usage = `
usage:

$ go-simple-proxy "<protocol>,<from_host:port>,<to_host:port>" ...
ex) go-simple-proxy "tcp,localhost:8080,example.com:80" "udp,localhost:8081,example2.com:80" "ws,localhost:8082,example3.com:80"
`

func init() {
	flag.BoolVar(&_verbose, "v", _verbose, "output verbose log")
	flag.BoolVar(&_debug, "vv", _debug, "output debug log")
	flag.StringVar(&_timeout, "timeout", _timeout, "set request/response timeout (e.g. 10s, 500ms, 1.5h)")
}

func main() {
	parseFlag()

	args, err := validateArg()
	if err != nil {
		exitWithError(err)
	}

	proxyList, err := NewProxtList(args)
	if err != nil {
		exitWithError(err)
	}
	proxyList.SetDefaultTimeout(timeout)

	loggingInfo("\n%s", proxyList.String())

	ch := make(chan os.Signal)
	signal.Notify(ch,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT)

	proxyList.ServeAll()
	loggingInfo("proxy started")

	s := <-ch
	switch s {
	case syscall.SIGINT:
		loggingInfo("\nkilled by SIGINT")
		proxyList.CloseAll()
	}
	os.Exit(0)
}

// parseFlag parses command line flag options.
func parseFlag() {
	flag.Parse()

	if _verbose || _debug {
		enableInfoLog()
		loggingInfo("enabled info log")
	}
	if _debug {
		enableDebugLog()
		loggingInfo("enabled debug log")
	}
	if _timeout != "" {
		var err error
		timeout, err = time.ParseDuration(_timeout)
		if err != nil {
			exitWithError(err)
		}
		loggingInfo("enabled timeout")
	}
}

// validateArg validates command line arguments.
func validateArg() (args []string, err error) {
	if len(os.Args) > 1 {
		for _, arg := range os.Args[1:] {
			if strings.HasPrefix(arg, "-") {
				continue
			}
			args = append(args, arg)
		}
	}

	if len(args) == 0 {
		return nil, fmt.Errorf("Argument for proxy is missing.\n%s", usage)
	}
	return args, nil
}

// exitWithError outputs error log and exits with error code.
func exitWithError(err error) {
	loggingError(err.Error())
	os.Exit(2)
}
