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

	if exists {
		s, ok := value.([]string)

		if !ok {
			fmt.Println("Something went wrong")
			conn.Write([]byte("+Something went wrong\r\n"))
		}
		for _, item := range s {
			slice = append(slice, item)
		}
	}
	for i := 1; i < len(args); i++ {
		slice = append(slice, args[i])
	}
	Data[key] = slice
	conn.Write(fmt.Appendf([]byte{}, ":%d\r\n", len(slice)))
}

func handleLRange(conn net.Conn, args []string) {
	key := args[0]
	low, err := strconv.ParseInt(args[1], 10, 64)

	if err != nil {
		fmt.Println("Cannot parse to int:", args[1])
		conn.Write(fmt.Appendf([]byte{}, "Cannot parse to int: %v\r\t", args[1]))
		return
	}
	high, err := strconv.ParseInt(args[2], 10, 64)
	if err != nil {
		fmt.Println("Cannot parse to int:", args[2])
		conn.Write(fmt.Appendf([]byte{}, "Cannot parse to int: %v\r\t", args[2]))
		return
	}

	value, exists := Data[key]

	if !exists {
		conn.Write([]byte("*0\r\n"))
		return
	}

	slice, ok := value.([]string)
	if !ok {
		conn.Write([]byte("+Invalid datatype\r\n"))
		return
	}

	if -low >= int64(len(slice)) {
		low = 0
	}
	if -high >= int64(len(slice)) {
		high = 0
	}

	if low < 0 {
		low = int64(len(slice)) + low
	}
	if high < 0 {
		high = int64(len(slice)) + high
	}

	if int(low) >= len(slice) || low > high {
		conn.Write([]byte("*0\r\n"))
		return
	}
	if int(high) >= len(slice) {
		high = int64(len(slice) - 1)
	}

	output := ""
	count := 0
	for i := low; i <= high; i++ {
		output = fmt.Sprintf("%s$%d\r\n%s\r\n", output, len(slice[i]), slice[i])
		count += 1
	}
	conn.Write([]byte(fmt.Sprintf("*%d\r\n%s", count, output)))
}

func handleLPush(conn net.Conn, args []string) {
	key := args[0]

	values, exists := Data[key]
	slice := []string{}
	if exists {
		s, ok := values.([]string)
		if !ok {
			fmt.Println("Unable to convert value to string slice")
			conn.Write([]byte("+Invalid datatype\r\n"))
			return
		}
		slice = s
	}

	for i := 1; i < len(args); i++ {
		slice = append([]string{args[i]}, slice...)
	}
	Data[key] = slice
	conn.Write(fmt.Appendf([]byte{}, ":%d\r\n", len(slice)))
}

func handleLLen(conn net.Conn, args []string) {
	key := args[0]
	content, exists := Data[key]
	if !exists {
		conn.Write([]byte(":0\r\n"))
		return
	}

	slice, ok := content.([]string)
	if !ok {
		fmt.Println("Unable to convert value to string slice")
		conn.Write([]byte("+Invalid datatype\r\n"))
		return
	}
	conn.Write(fmt.Appendf([]byte{}, ":%d\r\n", len(slice)))
}

func handleLPop(conn net.Conn, args []string) {
	key := args[0]
	content, exists := Data[key]
	if !exists {
		conn.Write([]byte("$-1\r\n"))
		return
	}

	slice, ok := content.([]string)
	if !ok {
		fmt.Println("Unable to convert value to string slice")
		conn.Write([]byte("+Invalid datatype\r\n"))
		return
	}
	if len(slice) == 0 {
		conn.Write([]byte("$-1\r\n"))
		return
	}

	toPop := slice[0]
	Data[key] = slice[1:]
	conn.Write(fmt.Appendf([]byte{}, "$%d\r\n%s\r\n", len(toPop), toPop))
}
