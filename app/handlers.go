package main

import (
	"fmt"
	"net"
	"strconv"
	"strings"
	"time"

	"github.com/xwhiz/redis-with-go/data"
)

func handleSetKey(conn net.Conn, args []string) {
	key, value := args[0], args[1]
	data.Data[key] = value
	conn.Write([]byte("+OK\r\n"))

	for index, arg := range args {
		if strings.ToLower(arg) != "px" {
			continue
		}

		sleepDuration, err := strconv.ParseInt(args[index+1], 10, 64)
		fmt.Printf("Sleeping %d for key: %v\n", sleepDuration, key)
		if err != nil {
			fmt.Println("Cannot parse to int:", args[index+1])
			continue
		}
		go func() {
			time.Sleep(time.Millisecond * time.Duration(sleepDuration))
			delete(data.Data, key)
		}()
	}
}

func handleGetKey(conn net.Conn, args []string) {
	key := args[0]
	fetched, exists := data.Data[key]

	if !exists {
		conn.Write([]byte("$-1\r\n"))
		return
	}

	value, ok := fetched.(string)
	if !ok {
		conn.Write([]byte("+Non-string type key\r\n"))
		return
	}

	conn.Write(fmt.Appendf([]byte{}, "$%d\r\n%s\r\n", len(value), value))
}

func handleRPush(conn net.Conn, args []string) {
	key := args[0]
	value, exists := data.Data[key]
	slice := []string{}

	if exists {
		s, ok := value.([]string)

		if !ok {
			fmt.Println("Something went wrong")
			conn.Write([]byte("+Something went wrong\r\n"))
		}
		slice = s
	}
	for i := 1; i < len(args); i++ {
		slice = append(slice, args[i])
	}
	conn.Write(fmt.Appendf([]byte{}, ":%d\r\n", len(args)-1))
}
