package main

import (
	"encoding/json"
	"fmt"
	"net"
	"slices"
	"sync"
)

type Message struct {
	Name string
	Msg  string
}

var (
	clients   = make([]net.Conn, 0)
	clientsMu sync.Mutex
)

func Find[T comparable](collection []T, el T) int {
	for i := range collection {
		if collection[i] == el {
			return i
		}
	}
	return -1
}

func removeClient(conn net.Conn) {
	clientsMu.Lock()
	defer clientsMu.Unlock()
	idx := Find(clients, conn)
	if idx > -1 {
		clients = slices.Delete(clients, idx, idx+1)
	}
}

func handleConnection(con net.Conn) {
	defer con.Close()
	clientAddr := con.RemoteAddr().String()
	fmt.Printf("Client connected: %s\n", clientAddr)
	d := json.NewDecoder(con)

	clientsMu.Lock()
	clients = append(clients, con)
	clientsMu.Unlock()
	firstCon := true
	for {
		var msg Message
		err := d.Decode(&msg)
		if err != nil {
			fmt.Printf("Client disconnected %s\n", clientAddr)
			removeClient(con)
			break
		}
		if firstCon {
			broadcastMsg(Message{Name: "[Server]", Msg: "New Player: " + string(msg.Name)}, con)
			firstCon = false
		} else {
			fmt.Printf("%s: %s\n", msg.Name, msg.Msg)
			broadcastMsg(msg, con)
		}
	}
}

func broadcastMsg(msg Message, sender net.Conn) {
	clientsMu.Lock()
	defer clientsMu.Unlock()
	for _, c := range clients {
		if c != sender {
			json.NewEncoder(c).Encode(msg)
		}
	}
}

func main() {
	listener, err := net.Listen("tcp", ":9000")
	if err != nil {
		panic(err)
	}
	defer listener.Close()

	for {
		con, err := listener.Accept()
		if err != nil {
			fmt.Printf("Error with client connection")
		}
		go handleConnection(con)
	}
}
