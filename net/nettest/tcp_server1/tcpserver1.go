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
	conn.SetOnExceptionCB(OnExceptionS)
	fmt.Printf("OnConnect: a new connection is ok, time at %s\n", time.Now().Format("2006-01-02 15:04:05"))
	gIndex++
	go jobS(conn)
}

func OnExceptionS(err error) {
	fmt.Printf("OnException: something wrong, time at %s\n", time.Now().Format("2006-01-02 15:04:05"))
}

func main() {
	tsvr, err := tcp.NewTcpServer("", 12345)
	if err != nil {
		fmt.Printf("main: NewTcpServer failed: %s\n", err.Error())
		return
	}
	err = tsvr.Start(OnConnectS)
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
		n, err := conn.ReadTimeout(b, 1500)
		// n, err := conn.Read(b)
		if err != nil {
			fmt.Printf("job: Read failed: %s\n", err.Error())
			break
		}
		if n == 0 {
			fmt.Println("read 0 bytes, continue reading...")
			time.Sleep(1 * time.Second)
			continue
		}
		fmt.Printf("job: receive message, size: %d, text: %s", n, string(b[:n]))

		n, err = buffw.WriteString(string(b[:n]))
		if err != nil {
			fmt.Printf("job: Write to file failed: %s\n", err.Error())
			break
		}
		buffw.Flush()

		rets := fmt.Sprintf("server receive message, size : %d", n)
		retb := []byte(rets)
		n, err = conn.Write(retb)
		if err != nil {
			fmt.Printf("job: Write failed: %s\n", err.Error())
			break
		}
		fmt.Printf("write to client ok, continue waiting for Read.\n")
	}

	fmt.Println("job: finish receive file, exit!")
}
