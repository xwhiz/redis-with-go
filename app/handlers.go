package main

import (
	"fmt"
	"net"
	"strconv"
	"strings"
	"time"
)

func handleSetKey(conn net.Conn, args []string) {
	key, value := args[0], args[1]
	Data[key] = value
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
			delete(Data, key)
		}()
	}
}

func handleGetKey(conn net.Conn, args []string) {
	key := args[0]
	fetched, exists := Data[key]

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
	value, exists := Data[key]
	slice := []string{}

	fmt.Println(Data, key, value, exists)

	if exists {
		s, ok := value.([]string)

		fmt.Println("this thing exists", s)

		if !ok {
			fmt.Println("Something went wrong")
			conn.Write([]byte("+Something went wrong\r\n"))
		}
		for _, item := range s {
			slice = append(slice, item)
		}
	}
	fmt.Println(slice)
	for i := 1; i < len(args); i++ {
		slice = append(slice, args[i])
	}
	fmt.Println(slice)
	Data[key] = slice
	conn.Write(fmt.Appendf([]byte{}, ":%d\r\n", len(slice)))
}
