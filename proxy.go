package main

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"net"
	"os"
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

	buf := &bytes.Buffer{}
	for {
		fmt.Printf("Start byte loop\n")
		data := make([]byte, 256)
		n, err := conn.Read(data)
		fmt.Printf("Done read\n")
		if err != nil {
			fmt.Printf("%v", err)
			return
		}
		buf.Write(data[:n])
		fmt.Printf("Done write\n")
		if data[0] == 13 && data[1] == 10 {
			break
		}
	}

	if _, err := rConn.Write(buf.Bytes()); err != nil {
		fmt.Printf("%v", err)
		return
	}
	log.Printf("sent:\n%v", hex.Dump(buf.Bytes()))

	data := make([]byte, 1024)
	n, err := rConn.Read(data)
	if err != nil {
		if err != io.EOF {
			fmt.Printf("%v", err)
			return
		} else {
			log.Printf("received err: %v", err)
		}
	}
	log.Printf("received:\n%v", hex.Dump(data[:n]))
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
