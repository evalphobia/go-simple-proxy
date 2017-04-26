package main

import (
	"fmt"
	"net"
)

// TCPProxy is proxy for TCP.
type TCPProxy struct {
	from string
	to   string

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
			loggingError("Connection on %s, %s", p.from, err.Error())
			l.Close()
			return
		}

		toReq, err := net.Dial("tcp", p.to)
		if err != nil {
			loggingError("Connection on %s, %s", p.to, err.Error())
			l.Close()
			return
		}
		proxyConn(fromReq, toReq)
	}
}
