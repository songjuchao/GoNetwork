package main

import (
	"os"
	"fmt"
	"net"
	"log"
)

func main() {
	if len(os.Args) != 3 {
		fmt.Fprintf(os.Stderr, "Usage: %s hostname port\n", os.Args[0])
		os.Exit(1)
	}

	hostname := os.Args[1]
	port := os.Args[2]
	serverTcpAddr, err := net.ResolveTCPAddr("tcp", hostname+":"+port)
	if err != nil {
		log.Fatal(err)
	}

	tcpListener, err := net.ListenTCP("tcp", serverTcpAddr)
	// NOTE: defer firstly or check err firstly
	defer tcpListener.Close()
	if err != nil {
		log.Fatal(err)
	}

	for {
		clientTcpConn, err := tcpListener.AcceptTCP()
		// NOTE
		defer clientTcpConn.Close()
		if err != nil {
			log.Fatal(err)
		}
		go echo(clientTcpConn)
	}
}

func echo(tcpConn *net.TCPConn) {
	buffer := make([]byte, 4096)
	for {
		size, err := tcpConn.Read(buffer)
		if err != nil {
			log.Println(err)
		}
		log.Printf("read size %d\n", size)

		size, err = tcpConn.Write(buffer[:size])
		if err != nil {
			log.Println(err)
		}
		log.Printf("write size %d\n", size)
	}
}
