package handler

import (
	"io"
	"log"
	"net"
	"net/http"
)

type connect struct {
}

func dialTarget(host string) (c net.Conn, err error) {
	raddr, err := net.ResolveTCPAddr("", host)
	if err != nil {
		return
	}

	c_, err := net.DialTCP("tcp", nil, raddr)
	if err != nil {
		return
	}

	c_.SetKeepAlive(true)
	c = c_

	return
}

func copy(dst io.ReadWriteCloser, src io.ReadWriteCloser) {
	io.Copy(dst, src)
	dst.Close()
	src.Close()
}

func (c *connect) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	if req.Method != "CONNECT" {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusMethodNotAllowed)
		io.WriteString(w, "405 must CONNECT\n")
		return
	}

	out, err := dialTarget(req.URL.Host)
	if err != nil {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusBadGateway)
		io.WriteString(w, "fail to connect target\n")
		return
	}

	in, _, err := w.(http.Hijacker).Hijack()
	if err != nil {
		out.Close()
		log.Print("connect hijacking ", req.RemoteAddr, ": ", err.Error())
		return
	}

	io.WriteString(in, "HTTP/1.0 200 Connected to target\r\n\r\n")
	go copy(in, out)
	go copy(out, in)
}

func New() http.Handler {
	return new(connect)
}