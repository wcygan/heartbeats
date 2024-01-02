package main

import (
	"fmt"
	"net"
	"os"
	"sync"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

const (
	port              = ":8080"
	heartbeatInterval = 5 * time.Second
)

var heartbeats = make(map[string]int)
var mutex = &sync.Mutex{}
var table = tview.NewTable()
var updateChan = make(chan peerStatus)

type peerStatus struct {
	peer   string
	status string
}

func main() {
	peers := os.Args[1:]

	go startServer()

	for _, ip := range peers {
		go sendHeartbeat(ip)
	}

	app := tview.NewApplication()
	table.SetBorder(true).SetTitle("Connections").SetTitleColor(tcell.ColorGreen)
	go func() {
		if err := app.SetRoot(table, true).Run(); err != nil {
			panic(err)
		}
	}()

	go func() {
		for ps := range updateChan {
			updateTable(ps.peer, ps.status)
		}
	}()

	select {}
}

func startServer() {
	ln, err := net.Listen("tcp", port)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer ln.Close()

	//fmt.Println("Server is running on port", port)

	for {
		conn, err := ln.Accept()
		if err != nil {
			//fmt.Println(err)
			continue
		}
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()

	buffer := make([]byte, 1024)
	peer := conn.RemoteAddr().String()
	for {
		_, err := conn.Read(buffer)
		if err != nil {
			fmt.Println(err)
			mutex.Lock()
			heartbeats[peer] = 0
			mutex.Unlock()
			updateChan <- peerStatus{peer, "(disconnected)"}
			return
		}
		//fmt.Println("Received:", string(buffer[:n]))
		mutex.Lock()
		heartbeats[peer]++
		mutex.Unlock()
		updateChan <- peerStatus{peer, fmt.Sprintf("%d", heartbeats[peer])}
	}
}

func sendHeartbeat(ip string) {
	for {
		//fmt.Printf("Attempting to connect to %s\n", ip)
		conn, err := net.Dial("tcp", ip+port)
		if err != nil {
			//fmt.Println(err)
			time.Sleep(heartbeatInterval)
			continue
		}
		defer conn.Close()

		for {
			_, err := conn.Write([]byte("Ping"))
			if err != nil {
				//fmt.Println(err)
				break
			}
			time.Sleep(heartbeatInterval)
		}
	}
}

func updateTable(peer, status string) {
	for row := 0; row < table.GetRowCount(); row++ {
		if cell := table.GetCell(row, 0); cell.Text == peer {
			table.SetCellSimple(row, 1, status)
			return
		}
	}
	table.SetCellSimple(table.GetRowCount(), 0, peer)
	table.SetCellSimple(table.GetRowCount()-1, 1, status)
}
