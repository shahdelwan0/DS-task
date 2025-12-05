package main

import (
	"bufio"
	"fmt"
	"net/rpc"
	"os"
	"strings"
)

type Message struct {
	UserID string
	Text   string
}

func main() {
	client, err := rpc.Dial("tcp", "localhost:1234")
	if err != nil {
		fmt.Println("Error connecting to server:", err)
		return
	}
	defer client.Close()

	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter your name: ")
	name, _ := reader.ReadString('\n')
	name = strings.TrimSpace(name)

	var ack string
	err = client.Call("ChatServer.Join", name, &ack)
	if err != nil {
		fmt.Println("Join error:", err)
		return
	}
	fmt.Println(ack)

	
	go func() {
		for {
			var msg string
			err := client.Call("ChatServer.Receive", name, &msg)
			if err == nil && msg != "" {
				fmt.Println("\n" + msg)
				fmt.Print("> ")
			}
		}
	}()

	fmt.Println("Connected to chatroom. Type messages below:")
	fmt.Println("Type 'exit' to leave.")

	
	for {
		fmt.Print("> ")
		text, _ := reader.ReadString('\n')
		text = strings.TrimSpace(text)
		if text == "exit" {
			fmt.Println("Exiting chat...")
			break
		}

		client.Call("ChatServer.SendMessage", Message{UserID: name, Text: text}, &ack)
	}
}
