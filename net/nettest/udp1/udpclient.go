package main

import (
	"bufio"
	"dadian/net/udp"
	"fmt"
	"os"
)

func main() {
	uc, err := udp.NewUdpClient("127.0.0.1", 12345, false)
	if err != nil {
		fmt.Printf("%s", err.Error())
		return
	}
	fmt.Printf("start to input:")
	input := bufio.NewScanner(os.Stdin)
	b := make([]byte, 1024*1024)
	for input.Scan() {
		s := input.Text()
		uc.Write([]byte(s))
		for {
			fmt.Printf("start to reading...\n")
			n, err := uc.ReadTimeout(b, 2000)
			if err != nil {
				fmt.Printf("read failed: %s\n", err.Error())
				break
			}
			if n == 0 {
				fmt.Printf("read 0 bytes from server, continue to input\n")
				fmt.Print("Input: ")
				break
			}
			fmt.Printf("read %d bytes from server: %s\n", n, b[:n])
			break
		}
	}
}
