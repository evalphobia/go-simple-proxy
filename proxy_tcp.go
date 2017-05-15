package main

import (
	"fmt"
	"net"
	"time"

	"gopkg.in/eapache/go-resiliency.v1/retrier"
)

// DefaultRetrier is default retry strategy when net.Dial to the destination is failed.
var DefaultRetrier = retrier.New(retrier.ExponentialBackoff(5, 200*time.Millisecond), nil)

// TCPProxy is proxy for TCP.
type TCPProxy struct {
	from    string
	to      string
	timeout time.Duration

	listner net.Listener // local port server
	retry   *retrier.Retrier
}

// newTCPProxy creates TCPProxy with initialized local port listener.
func newTCPProxy(from, to string) (*TCPProxy, error) {
	l, err := net.Listen("tcp", from)
	if err != nil {
		return nil, err
	}
	return &TCPProxy{
		from:    from,
		to:      to,
		listner: l,
		retry:   DefaultRetrier,
	}, nil
}

func (p TCPProxy) String() string {
	return fmt.Sprintf("TCPProxy: %s -> %s", p.from, p.to)
}

// Close closes local port listner.
func (p *TCPProxy) Close() {
	if p.listner != nil {
		p.listner.Close()
	}
}

// Serve serves proxy network.
func (p *TCPProxy) Serve() {
	l := p.listner
	for {
		fromReq, err := l.Accept()
		if err != nil {
			loggingError("Connection from %s, %s", p.from, err.Error())
			continue
		}

		err = p.retry.Run(func() error {
			toReq, err := net.Dial("tcp", p.to)
			if err != nil {
				loggingError("Connection to %s, %s", p.to, err.Error())
				return err
			}

			TCPPipe{
				From:    fromReq,
				To:      toReq,
				Timeout: p.timeout,
				Debug:   _debug,
			}.Do()
			return nil
		})
		if err != nil {
			loggingError("Give up connection to %s, %s", p.to, err.Error())
		}
	}
}

func (p *TCPProxy) SetTimeout(t time.Duration) {
	p.timeout = t
}
