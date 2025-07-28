package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

var Data = map[string]any{}

func main() {
	fmt.Println("Logs from your program will appear here!")

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
			handleSetKey(conn, args)
		}
		if command == "GET" {
			handleGetKey(conn, args)
		}
		if command == "RPUSH" {
			handleRPush(conn, args)
		}
		if command == "LRANGE" {
			handleLRange(conn, args)
		}
		if command == "LPUSH" {
			handleLPush(conn, args)
		}
	}
}
