package main

import (
	"fmt"
	"net"
	"time"
)

// TCPProxy is proxy for TCP.
type TCPProxy struct {
	from    string
	to      string
	timeout time.Duration

	listner net.Listener // local port server
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

		toReq, err := net.Dial("tcp", p.to)
		if err != nil {
			loggingError("Connection to %s, %s", p.to, err.Error())
			continue
		}

		TCPPipe{
			From:    fromReq,
			To:      toReq,
			Timeout: p.timeout,
			Debug:   _debug,
		}.Do()
	}
}

func (p *TCPProxy) SetTimeout(t time.Duration) {
	p.timeout = t
}
