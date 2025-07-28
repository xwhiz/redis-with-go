package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

var data = map[string]string{}

func main() {
	// You can use print statements as follows for debugging, they'll be visible when running tests.
	fmt.Println("Logs from your program will appear here!")

	// Uncomment this block to pass the first stage

	l, err := net.Listen("tcp", "0.0.0.0:6379")
	if err != nil {
		fmt.Println("Failed to bind to port 6379")
		os.Exit(1)
	}
	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting connection: ", err.Error())
			continue
		}

		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	scanner := bufio.NewScanner(conn)
	defer conn.Close()

	for scanner.Scan() {
		command, args := parseCommands(scanner)

		if command == "PING" {
			conn.Write([]byte("+PONG\r\n"))
		}
		if command == "ECHO" {
			conn.Write(fmt.Appendf([]byte{}, "$%d\r\n%s\r\n", len(args[0]), args[0]))
		}
		if command == "SET" {
			key, value := args[0], args[1]
			data[key] = value
			conn.Write([]byte("+OK\r\n"))
		}
		if command == "GET" {
			key := args[0]
			value, exists := data[key]

			if exists {
				conn.Write(fmt.Appendf([]byte{}, "$%d\r\n%s\r\n", len(value), value))
			} else {
				conn.Write([]byte("$-1\r\n"))
			}
		}
	}
}
