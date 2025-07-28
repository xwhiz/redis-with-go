package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
)

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

	tokens := []string{}
	var argsCount int64 = 0
	for scanner.Scan() {
		text := scanner.Text()

		if strings.HasPrefix(text, "*") {
			i, err := strconv.ParseInt(strings.TrimPrefix(text, "*"), 10, 64)
			if err != nil {
				fmt.Printf("Unable to convert string to int: %v\n", err)
			}
			argsCount = i
		}

		for range argsCount {
			scanner.Scan()
			scanner.Scan()
			text := scanner.Text()
			tokens = append(tokens, text)
		}

		command := tokens[0]

		if command == "PING" {
			conn.Write([]byte("+PONG\r\n"))
		}
		if command == "ECHO" {
			conn.Write(fmt.Appendf([]byte{}, "$%d\r\n%s\r\n", len(tokens[1]), tokens[1]))
		}

		tokens = []string{}
	}
}
