package main

import (
	"io"
	"net"
	"sync/atomic"
	"time"
)

type TCPPipe struct {
	From    net.Conn
	To      net.Conn
	Timeout time.Duration
	Debug   bool

	status int32
	timer  *time.Timer
}

func (p TCPPipe) Do() {
	if p.Debug {
		p.doDebug()
		return
	}
	p.do()
}

func (p *TCPPipe) do() {
	p.closeAfterTimeout()

	// response
	go func() {
		if err := pipe(p.To, p.From); err != nil {
			loggingError("Proxy response on %s:%s, %s", p.From.RemoteAddr(), p.To.LocalAddr(), err.Error())
		}
		p.doneTo()
	}()

	// request
	go func() {
		if err := pipe(p.From, p.To); err != nil {
			loggingError("Proxy request on %s:%s, %s", p.From.LocalAddr(), p.To.RemoteAddr(), err.Error())
		}
		p.doneFrom()
	}()
}

func (p *TCPPipe) doDebug() {
	p.closeAfterTimeout()
	isRequest := true

	// response
	go func() {
		if err := pipeDebug(p.To, p.From, !isRequest); err != nil {
			loggingError("Proxy response on %s:%s, %s", p.From.RemoteAddr(), p.To.LocalAddr(), err.Error())
		}
		p.doneTo()
	}()

	// request
	go func() {
		if err := pipeDebug(p.From, p.To, isRequest); err != nil {
			loggingError("Proxy request on %s:%s, %s", p.From.LocalAddr(), p.To.RemoteAddr(), err.Error())
		}
		p.doneFrom()
	}()
}

func (p *TCPPipe) closeAfterTimeout() {
	if int(p.Timeout) < 1 {
		return
	}

	p.timer = time.AfterFunc(p.Timeout, func() {
		p.To.Close()
		p.From.Close()
	})
}

func (p *TCPPipe) doneFrom() {
	if atomic.AddInt32(&p.status, -1) == 0 {
		p.stopTimer()
	}
}

func (p *TCPPipe) doneTo() {
	if atomic.AddInt32(&p.status, 1) == 0 {
		p.stopTimer()
	}
}

func (p *TCPPipe) stopTimer() {
	if p.timer != nil {
		p.timer.Stop()
	}
}

// pipe sends and receives the network stream.
func pipe(from, to net.Conn) error {
	_, err := io.Copy(to, from)
	return err
}

// pipeDebug sends and receives the network stream.
func pipeDebug(from, to net.Conn, isRequest bool) error {
	var direction string
	var fromAddr, toAddr net.Addr
	if isRequest {
		direction = "request"
		fromAddr = from.LocalAddr()
		toAddr = to.RemoteAddr()
	} else {
		direction = "response"
		fromAddr = from.RemoteAddr()
		toAddr = to.LocalAddr()
	}

	buf := make([]byte, 0, 4096)
	tmp := make([]byte, 256)
	for {
		n, err := from.Read(tmp)
		if err != nil {
			if err != io.EOF {
				loggingError("[pipeDebug %s: %s -> %s] %s", direction, fromAddr, toAddr, err.Error())
			}
			break
		}
		_, err = to.Write(tmp)
		if err != nil {
			return err
		}
		buf = append(buf, tmp[:n]...)
	}

	loggingDebug("[Pipe %s: %s -> %s] %s", direction, fromAddr, toAddr, buf)
	return nil
}
