package main

import (
	"log"
)

func main() {
	log.Println("Welcome to http2socks.")
	dialer := NewSocksDialer("tcp", "127.0.0.1:1080")
	listener := NewListener("tcp", "127.0.0.1", &HTTPHandler{Dialer: dialer})

	listener.Serve()

	quit := make(chan bool)
	<-quit
}
