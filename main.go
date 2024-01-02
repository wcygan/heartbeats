package main

import (
	"fmt"
	"net"
	"os"
	"time"
)

const (
	port              = ":8080"
	heartbeatInterval = 5 * time.Second
)

func main() {
	ips := os.Args[1:]

	go startServer()

	for _, ip := range ips {
		go sendHeartbeat(ip)
	}

	select {}
}

func startServer() {
	ln, err := net.Listen("tcp", port)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer ln.Close()

	fmt.Println("Server is running on port", port)

	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println(err)
			continue
		}
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()

	buffer := make([]byte, 1024)
	for {
		n, err := conn.Read(buffer)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println("Received:", string(buffer[:n]))
	}
}

func sendHeartbeat(ip string) {
	for {
		fmt.Printf("Attempting to connect to %s\n", ip)
		conn, err := net.Dial("tcp", ip+port)
		if err != nil {
			fmt.Println(err)
			time.Sleep(heartbeatInterval)
		}
		defer conn.Close()

		for {
			_, err := conn.Write([]byte("Hello"))
			if err != nil {
				fmt.Println(err)
				return
			}
			time.Sleep(heartbeatInterval)
		}
	}
}
