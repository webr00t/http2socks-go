package main

import (
	"bufio"
	"log"
	"net"
	"net/http"
)

const (
	connectionEstablished = "200 Connection Established"
)

// HTTPHandler is the http implementation of Handler.
type HTTPHandler struct {
	Dialer Dialer
}

// Handle responses http tunnel request.
func (h *HTTPHandler) Handle(conn net.Conn) {
	defer conn.Close()
	req, err := http.ReadRequest(bufio.NewReader(conn))
	if err != nil {
		//log.Err(err).Msgf("[http] %s -> %s: valid request", conn.RemoteAddr(), conn.LocalAddr())
		log.Printf("[http] %s -> %s: valid request", conn.RemoteAddr(), conn.LocalAddr())
		return
	}
	defer req.Body.Close()

	h.handleRequest(conn, req)
}

func (h *HTTPHandler) handleRequest(conn net.Conn, req *http.Request) {
	host := req.Host
	if _, port, _ := net.SplitHostPort(host); port == "" {
		host = net.JoinHostPort(host, "80")
	}

	resp := &http.Response{
		ProtoMajor: 1,
		ProtoMinor: 1,
		Header:     http.Header{},
	}

	if req.Method != http.MethodConnect && req.URL.Scheme != "http" {
		resp.StatusCode = http.StatusBadRequest
		resp.Write(conn)
		return
	}

	cc, err := h.Dialer.Dial("tcp", host)
	if err != nil {
		resp.StatusCode = http.StatusServiceUnavailable
		//log.Err(err).Msgf("[http] %s -> %s ->%s: tcp connect failed", conn.RemoteAddr(), conn.LocalAddr(), host)
		log.Printf("[http] %s -> %s ->%s: tcp connect failed", conn.RemoteAddr(), conn.LocalAddr(), host)
		resp.Write(conn)
		return
	}
	defer cc.Close()

	resp.StatusCode = http.StatusOK
	resp.Status = connectionEstablished
	resp.Header = http.Header{}

	resp.Write(conn)
	//log.Info().Msgf("[http] %s -> %s -> %s: success", conn.RemoteAddr(), conn.LocalAddr(), host)
	log.Printf("[http] %s -> %s -> %s: success", conn.RemoteAddr(), conn.LocalAddr(), host)
	transport(conn, cc)
	//log.Info().Msgf("[http] %s -> %s -> %s: closed", conn.RemoteAddr(), conn.LocalAddr(), host)
	log.Printf("[http] %s -> %s -> %s: closed", conn.RemoteAddr(), conn.LocalAddr(), host)
}
