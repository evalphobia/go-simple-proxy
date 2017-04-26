package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

var (
	verbose bool
)

var usage = `
usage:

$ go-simple-proxy "<protocol>,<from_host:port>,<to_host:port>" ...
ex) go-simple-proxy "tcp,localhost:8080,example.com:80" "udp,localhost:8081,example2.com:80" "ws,localhost:8082,example3.com:80"
`

func init() {
	flag.BoolVar(&verbose, "v", verbose, "output verbose log")
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
	if !verbose {
		// TODO: implements flag options
		// disableLog()
	}
}

// validateArg validates command line arguments.
func validateArg() (args []string, err error) {
	if len(os.Args) < 2 {
		return nil, fmt.Errorf("Argument for proxy is missing.\n%s", usage)
	}
	return os.Args[1:], nil
}

// exitWithError outputs error log and exits with error code.
func exitWithError(err error) {
	loggingError(err.Error())
	os.Exit(2)
}
