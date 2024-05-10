package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
)

func main() {
	// Connect to this socket
	//It is used for establishing connections to various network services.
	// dial  It is used to establish a connection to a remote network address over TCP
	conn, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		fmt.Println("Error connecting:", err.Error())
		return
	}

	for {
		// Read input from stdin
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("Enter command: ")
		text, _ := reader.ReadString('\n')

		// Send text to conn object
		fmt.Fprintf(conn, text+"\n")

		// Listen for reply and read his response until send \n
		message, _ := bufio.NewReader(conn).ReadString('\n')
		message = strings.TrimSuffix(message, "\n")
		messageAll := strings.Split(message, "###")
		if len(messageAll) == 1 {
			fmt.Println("Message from server: " + strings.TrimSpace(messageAll[0]))
		}
		for i := 0; i < len(messageAll)-1; i++ {
			fmt.Println("Message from server: " + strings.TrimSpace(messageAll[i]))
		}
	}
}
