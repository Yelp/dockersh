package main

import (
	"fmt"
	"io"
	"net"
)

func proxyConn(remoteAddr string, conn *net.TCPConn) {
	rAddr, err := net.ResolveTCPAddr("tcp", remoteAddr)
	if err != nil {
		fmt.Printf("%v", err)
		return
	}

	rConn, err := net.DialTCP("tcp", nil, rAddr)
	if err != nil {
		fmt.Printf("%v", err)
		return
	}
	defer rConn.Close()

	go io.Copy(conn, rConn)
	io.Copy(rConn, conn)
}

func handleConn(remoteAddr string, in <-chan *net.TCPConn, out chan<- *net.TCPConn) {
	for conn := range in {
		proxyConn(remoteAddr, conn)
		out <- conn
	}
}

func closeConn(in <-chan *net.TCPConn) {
	for conn := range in {
		conn.Close()
	}
}

func proxyMain(localAddr string, remoteAddr string) {
	fmt.Printf("Listening: %v\nProxying: %v\n\n", localAddr, remoteAddr)

	addr, err := net.ResolveTCPAddr("tcp", localAddr)
	if err != nil {
		fmt.Printf("%v", err)
		return
	}

	listener, err := net.ListenTCP("tcp", addr)
	if err != nil {
		fmt.Printf("%v", err)
		return
	}

	pending, complete := make(chan *net.TCPConn), make(chan *net.TCPConn)

	for i := 0; i < 5; i++ {
		go handleConn(remoteAddr, pending, complete)
	}
	go closeConn(complete)

	for {
		conn, err := listener.AcceptTCP()
		if err != nil {
			fmt.Printf("%v", err)
			return
		}
		pending <- conn
	}
}
