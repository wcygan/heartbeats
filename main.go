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

type PeerStatus struct {
	peer   string
	action Action
}

type Action int

const (
	CONNECT Action = iota
	DISCONNECT
	HEARTBEAT
)

func main() {
	peers := os.Args[1:]

	go startServer()

	for _, ip := range peers {
		go sendHeartbeat(ip)
	}

	select {}
}

func startServer() {
	listener, err := net.Listen("tcp", port)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer listener.Close()

	var updates = make(chan PeerStatus)
	go updateTerminalUI(updates)

	for {
		conn, err := listener.Accept()
		if err != nil {
			continue
		}
		go handleConnection(conn, updates)
	}
}

func updateTerminalUI(updates <-chan PeerStatus) {
	for update := range updates {
		switch update.action {
		case CONNECT:
			fmt.Printf("%s connected\n", update.peer)
		case DISCONNECT:
			fmt.Printf("%s disconnected\n", update.peer)
		case HEARTBEAT:
			fmt.Printf("%s heartbeat\n", update.peer)
		}
	}
}

func handleConnection(conn net.Conn, updates chan<- PeerStatus) {
	defer conn.Close()

	updates <- PeerStatus{conn.RemoteAddr().String(), CONNECT}

	buffer := make([]byte, 1024)
	peer := conn.RemoteAddr().String()
	for {
		_, err := conn.Read(buffer)
		if err != nil {
			updates <- PeerStatus{peer, DISCONNECT}
			return
		}
		updates <- PeerStatus{peer, HEARTBEAT}
	}
}

func sendHeartbeat(ip string) {
	for {
		fmt.Printf("Sending heartbeat to %s\n", ip)
		conn, err := net.Dial("tcp", ip+port)
		if err != nil {
			fmt.Printf("Could not connect to %s\n", ip)
			time.Sleep(heartbeatInterval)
			continue
		}
		defer conn.Close()

		for {
			fmt.Printf("Sending heartbeat to %s\n", ip)
			_, err := conn.Write([]byte("Ping"))
			if err != nil {
				fmt.Printf("Could not send heartbeat to %s\n", ip)
				break
			}
			time.Sleep(heartbeatInterval)
		}
	}
}
