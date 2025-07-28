package main

import "fmt"

func getRespString(slice []string, low int, high int) []byte {
	output := ""
	count := 0
	for i := low; i <= high; i++ {
		output = fmt.Sprintf("%s$%d\r\n%s\r\n", output, len(slice[i]), slice[i])
		count += 1
	}
	return fmt.Appendf([]byte{}, "*%d\r\n%s", count, output)
}
