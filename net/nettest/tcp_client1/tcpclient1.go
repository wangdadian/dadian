package main

import (
	"bufio"
	"dadian/net/tcp"
	"fmt"
	"io"
	"os"
	"time"
)

func OnConnect(conn *tcp.TcpConn) {
	fmt.Printf("OnConnect: client connect ok, time at %s\n", time.Now().Format("2006-01-02 15:04:05"))
	go job(conn)
}

func OnException(errType int) {
	fmt.Printf("OnException: something wrong, time at %s\n", time.Now().Format("2006-01-02 15:04:05"))
}

func main() {
	tclient, err := tcp.NewTcpClient("192.168.1.119", 12345)
	if err != nil {
		fmt.Printf("main: NewTcpClient failed: %s\n", err.Error())
		return
	}
	err = tclient.Connect(true, OnConnect, OnException)
	if err != nil {
		fmt.Printf("main: Connect failed: %s\n", err.Error())
		return
	}
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
			break
		}
		fmt.Printf("%s job: read bytes size: %d, text:\n", time.Now().Format("2006-01-02 15:04:05"), len(b))
		s := string(b)
		fmt.Println(s)
		n, err := conn.Write(b)
		if err != nil {
			fmt.Printf("job: send to server failed: %s\n", err.Error())
			return
		}

		fmt.Println("job: send to server, len: ", n)
		n, err = conn.Read(rb)
		if err != nil {
			fmt.Printf("job: read form server failed: %s\n", err.Error())
			break
		}
		fmt.Printf("%s receive message from server, size: %d, text:\n", time.Now().Format("2006-01-02 15:04:05"), n)
		s = string(rb[:n])
		fmt.Println(s)
		time.Sleep(500 * time.Millisecond)
	}

	fmt.Println("job: finish send file, exit!")
}
