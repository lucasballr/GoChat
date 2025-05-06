package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
	"os"
	"strings"
)

type Client struct {
	Name string
	Msg  string
}

func printMessage(client Client, num int) {
	if num == 0 {
		fmt.Printf("\r")
	} else {
		fmt.Printf("\033[%dA", num)
	}
	fmt.Printf("%s: %s\n", client.Name, client.Msg)
	fmt.Printf("> ")
}

func main() {
	serverAddr := "localhost:9000"
	name := string(os.Args[1])
	reader := bufio.NewReader(os.Stdin)
	con, err := net.Dial("tcp", serverAddr)
	if err != nil {
		panic(err)
	}
	defer con.Close()

	go func() {
		d := json.NewDecoder(con)
		for {
			var msg Client
			err := d.Decode(&msg)
			if err != nil {
				fmt.Println("Server Error")
				os.Exit(0)
			}
			printMessage(msg, 0)
		}
	}()

	json.NewEncoder(con).Encode(Client{Name: name, Msg: "INITIAL CONNECTION"})
	fmt.Printf("> ")
	for {
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)
		if input == "q" {
			fmt.Println("Goodbye.")
			os.Exit(0)
		}
		client := Client{Name: name, Msg: input}
		json.NewEncoder(con).Encode(client)
		printMessage(client, 1)
	}
}
