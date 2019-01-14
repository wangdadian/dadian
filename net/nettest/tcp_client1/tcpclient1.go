package main

import (
	"bufio"
	"dadian/net/tcp"
	"fmt"
	"io"
	"os"
	"time"
)

func OnException(err error) {
	fmt.Printf("OnException: something wrong, time at %s\n", time.Now().Format("2006-01-02 15:04:05"))
}

func main() {
	tcpconn, err := tcp.Connect("127.0.0.1", 12345, true)
	if err != nil {
		fmt.Printf("main: Connect failed: %s\n", err.Error())
		return
	}
	tcpconn.SetOnExceptionCB(OnException)
	go job(tcpconn)
	for {
		time.Sleep(500 * time.Millisecond)
	}
}

func job(conn *tcp.TcpConn) {
	file, err := os.Open("./log.txt")
	if err != nil {
		fmt.Printf("job: open file failed: %s\n", err.Error())
		return
	}
	defer file.Close()
	defer conn.Close()
	rb := make([]byte, 1024)
	buffr := bufio.NewReader(file)
	for {
		b, err := buffr.ReadBytes('\n')
		if err == io.EOF {
			file.Seek(0, 0)
			buffr = bufio.NewReader(file)
			continue
		}
		// fmt.Printf("job: read bytes size: %d, text:\n", len(b))
		s := string(b)
		fmt.Println(s)
		n, err := conn.Write(b)
		if err != nil {
			fmt.Printf("job: send to server failed: %s\n", err.Error())
			return
		}

		fmt.Println("job: send to server, len: ", n)
		n, err = conn.ReadTimeout(rb, 1500)
		if n == 0 && err != nil {
			fmt.Printf("job: read form server failed: %s\n", err.Error())
			break
		}
		fmt.Printf("job: read from server ok, message: %s\n", string(rb[:n]))
		time.Sleep(500 * time.Millisecond)
	}

	fmt.Println("job: finish send file, exit!")
}
