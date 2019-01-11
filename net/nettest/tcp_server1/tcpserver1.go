package main

import (
	"bufio"
	"dadian/net/tcp"
	"fmt"
	"os"
	"time"
)

var gIndex int

func init() {
	gIndex = 0
}
func OnConnectS(conn *tcp.TcpConn) {
	fmt.Printf("OnConnect: a new connection is ok, time at %s\n", time.Now().Format("2006-01-02 15:04:05"))
	gIndex++
	go jobS(conn)
}

func OnExceptionS(errType int) {
	fmt.Printf("OnException: something wrong, time at %s\n", time.Now().Format("2006-01-02 15:04:05"))
}

func main() {
	tsvr, err := tcp.NewTcpServer("", 12345)
	if err != nil {
		fmt.Printf("main: NewTcpServer failed: %s\n", err.Error())
		return
	}
	err = tsvr.Start(OnConnectS, OnExceptionS)
	if err != nil {
		fmt.Printf("main: Start failed: %s\n", err.Error())
		return
	}
	for {
		time.Sleep(500 * time.Millisecond)
	}
}

func jobS(conn *tcp.TcpConn) {
	szFile := fmt.Sprintf("./log-svr%d.txt", gIndex)
	file, err := os.OpenFile(szFile, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0)
	if err != nil {
		fmt.Printf("job: open|create file failed: %s\n", err.Error())
		return
	}
	defer file.Close()
	defer conn.Close()
	buffw := bufio.NewWriter(file)
	b := make([]byte, 1024)
	for {
		n, err := conn.Read(b)
		if err != nil {
			fmt.Printf("job: Read failed: %s\n", err.Error())
			break
		}
		fmt.Printf("%s job: receive message, size: %d, text:\n", time.Now().Format("2006-01-02 15:04:05"), n)
		s := string(b[:n])
		fmt.Println(s)

		// n, err = buffw.Write(b)
		n, err = buffw.WriteString(s)
		if err != nil {
			fmt.Printf("job: Write failed: %s\n", err.Error())
			break
		}
		fmt.Printf("job: Write ok, size: %d\n", n)
		buffw.Flush()

		rets := fmt.Sprintf("[%s] server receive message, size : %d, text: \n", time.Now().Format("2006-01-02 15:04:05"), n)
		retb := []byte(rets)
		n, err = conn.Write(retb)
		if err != nil {
			fmt.Printf("job: Write failed: %s\n", err.Error())
			break
		}
		fmt.Printf("%s continue waiting for Read.\n", time.Now().Format("2006-01-02 15:04:05"))
	}

	fmt.Println("job: finish receive file, exit!")
}
