package main

import (
	"bufio"
	"fmt"
	"strconv"
	"strings"
)

func parseCommands(scanner *bufio.Scanner) (string, []string) {
	argsCount := int64(0)
	text := scanner.Text()

	if strings.HasPrefix(text, "*") {
		i, err := strconv.ParseInt(strings.TrimPrefix(text, "*"), 10, 64)
		if err != nil {
			fmt.Printf("Unable to convert string to int: %v\n", err)
		}
		argsCount = i
	}

	tokens := []string{}

	for range argsCount {
		scanner.Scan()
		scanner.Scan()
		text := scanner.Text()
		tokens = append(tokens, text)
	}

	return tokens[0], tokens[1:]
}
