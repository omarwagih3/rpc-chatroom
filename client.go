package main

import (
	"bufio"
	"fmt"
	"log"
	"net/rpc"
	"os"
	"strings"
)

type Message struct {
	Sender  string
	Content string
}

func main() {
	client, err := rpc.Dial("tcp", "localhost:1234")
	if err != nil {
		log.Fatalf("Dialing error: %v", err)
	}
	defer client.Close()

	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Enter your name: ")
	name, _ := reader.ReadString('\n')
	name = strings.TrimSpace(name)

	fmt.Println("Welcome to the chatroom! Type 'exit' to quit.")

	for {
		fmt.Print("Enter message: ")
		text, _ := reader.ReadString('\n')
		text = strings.TrimSpace(text)

		if text == "exit" {
			fmt.Println("Goodbye!")
			break
		}

		var messages []Message
		msg := Message{Sender: name, Content: text}

		err = client.Call("ChatServer.SendMessage", msg, &messages)
		if err != nil {
			log.Printf("Error sending message: %v", err)
			break
		}

		fmt.Println("\n--- Chat History ---")
		for _, m := range messages {
			fmt.Printf("%s: %s\n", m.Sender, m.Content)
		}
		fmt.Println("--------------------\n")
	}
}
