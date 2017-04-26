package main

import (
	"fmt"
	"io"
	"net"
	"strings"
)

// ProxyList contains multiple Proxy and handle them.
type ProxyList struct {
	list []Proxy
}

// NewProxtList returns initialized ProxyList by argument list.
// Each argument must be follow the format "<protocol>,<from server>:<from port>,<to server>:<to port>".
func NewProxtList(args []string) (*ProxyList, error) {
	proxyList := make([]Proxy, len(args))
	for i, arg := range args {
		var err error
		proxyList[i], err = CreateProxy(arg)
		if err != nil {
			return nil, err
		}
	}

	return &ProxyList{
		list: proxyList,
	}, nil
}

func (l ProxyList) String() string {
	result := make([]string, len(l.list))
	for i, p := range l.list {
		result[i] = p.String()
	}
	return strings.Join(result, "\n")
}

// ServeAll serves all proxy network.
func (l *ProxyList) ServeAll() {
	for _, p := range l.list {
		go p.Serve()
	}
}

// CloseAll closes all proxy network.
func (l *ProxyList) CloseAll() {
	for _, p := range l.list {
		p.Close()
	}
}

// Proxy is type of proxy network.
type Proxy interface {
	Serve()
	Close()
	String() string
}

// CreateProxy is a factory function for Proxy.
// Choose correct type from the given argument.
func CreateProxy(setting string) (Proxy, error) {
	part := strings.Split(setting, ",")
	if len(part) != 3 {
		return nil, fmt.Errorf("cannot parse arg, '%s'\n%s", setting, usage)
	}

	switch part[0] {
	case "tcp":
		return newTCPProxy(part[1], part[2])
	// case "udp":
	// return newProxyUDP(part[1], part[2]), nil
	// case "ws", "websocket":
	// return newProxyWebSocket(part[1], part[2]), nil
	default:
		return nil, fmt.Errorf("unknown protocol, '%s' from '%s'\n%s", part[0], setting, usage)
	}
}

// proxyConn sends and receives the network stream.
func proxyConn(from, to net.Conn) {
	go func() {
		if _, err := io.Copy(to, from); err != nil {
			loggingError("Proxy request on %s:%s, %s", from.LocalAddr(), to.LocalAddr(), err.Error())
		}
	}()
	go func() {
		if _, err := io.Copy(from, to); err != nil {
			loggingError("Proxy response on %s:%s, %s", from.LocalAddr(), to.LocalAddr(), err.Error())
		}
	}()
}
