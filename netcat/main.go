package main

import (
	"flag"
	"net"
	"os"
	"fmt"
	"log"
)

var l = flag.Bool("l", false, "Used to specify that nc should listen for an incoming connection")

// FIXME: don't support concurrent connection
func main() {
	flag.Parse()

	if len(os.Args) < 3 {
		fmt.Printf("Usage:\n %s hostname port\n %s -l port\n", os.Args[0], os.Args[0])
		os.Exit(1)
	}

	if *l {
		// FIXME
		// NOTE: the correct way to create listen TCP
		port := os.Args[2]
		hostTcpAddr, err := net.ResolveTCPAddr("tcp","127.0.0.1:"+port)
		if err != nil {
			log.Fatal(err)
		}

		tcpListener, err := net.ListenTCP("tcp", hostTcpAddr)
		if err != nil {
			log.Fatal(err)
		}

		for {
			clientTcpConn, err := tcpListener.AcceptTCP()
			// defer clientTcpConn.Close()
			if err != nil {
				log.Print(err)
				continue
			}
			netcat(clientTcpConn)
		}
	} else {
		hostname := os.Args[1]
		port := os.Args[2]
		clientTcpAddr, err := net.ResolveTCPAddr("tcp", hostname+":"+port)
		if err != nil {
			log.Fatal(err)
		}

		tcpConn, err := net.DialTCP("tcp", nil, clientTcpAddr)
		// defer tcpConn.Close()
		if err != nil {
			log.Fatal(err)
		}
		netcat(tcpConn)
	}
}

func netcat(conn *net.TCPConn) {
	// read from stdin, write to socket
	go func() {
		buffer := make([]byte, 4096)
		for {
			size, err := os.Stdin.Read(buffer)
			if err != nil {
				log.Println(err)
			}

			if size != 0 {
				conn.Write(buffer[:size])
			} else {
				break;
			}
		}

		// 从标准输入读数据返回0 没有数据可写 关闭连接的写入端 保证已写如的数据仍可读
		conn.CloseWrite()
	}()

	buffer := make([]byte, 4096)
	for {
		size, err := conn.Read(buffer)
		if err != nil {
			log.Println(err)
		}

		if size != 0 {
			os.Stdout.Write(buffer)
		} else {
			break
		}
	}

	// 从conn读数据返回0 没有数据可读 直接关闭conn
	conn.Close()
}