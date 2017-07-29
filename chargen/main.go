package main

import (
	"os"
	"fmt"
	"net"
	"log"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Printf("Usage:\n %s hostname port\n", os.Args[0])
		os.Exit(1)
	}

	hostname := os.Args[1]
	port := os.Args[2]
	svrTcpAddr, err := net.ResolveTCPAddr("tcp", hostname+":"+port)
	if err != nil {
		log.Fatal(err)
	}

	tcpListener, err := net.ListenTCP("tcp", svrTcpAddr)
	if err != nil {
		log.Fatal(err)
	}

	for {
		clientTcpConn, err := tcpListener.AcceptTCP()
		if err != nil {
			log.Print(err)
			continue
		}
		go chargen(clientTcpConn)
	}
}

func chargen(conn *net.TCPConn) {
	var line []byte
	for i := 33; i < 127; i++ {
		line = append(line, byte(i))
	}
	line = append(line, line...)

	var msg string
	for i := 0; i < 127-33; i++ {
		msg += string(line[i:i+72])+"\n"
	}

	for {
		conn.Write([]byte(msg))
	}
}