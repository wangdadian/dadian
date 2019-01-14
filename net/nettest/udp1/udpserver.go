package main

import (
	"dadian/net/udp"
	"fmt"
	"time"
)

func main() {
	uc, err := udp.NewUdpServer("127.0.0.1", 12345)
	if err != nil {
		fmt.Printf("%s", err.Error())
		return
	}
	byASCII := []byte("ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	iIndex := 0
	iLen := len(byASCII)
	for {
		fmt.Printf("start to read...\n")
		b := make([]byte, 50*1024)
		n, addr, err := uc.ReadFromUdpAddrTM(b, 2000)
		if err != nil {
			fmt.Printf("read failed: %s\n", err.Error())
			return
		}
		if n == 0 {
			fmt.Printf("read 0 bytes, continue reading...\n")
			time.Sleep(1 * time.Second)
			continue
		}
		fmt.Printf("read %d bytes from [%s], text: %s\n", n, addr.String(), string(b[:n]))
		//n, err = uc.WriteToUdpAddr([]byte(fmt.Sprintf("## server recv %d bytes", n)), addr)
		bySend := make([]byte, 65507) // 最多一次发送65507个字节，多了会报错
		if iIndex >= iLen {
			iIndex = 0
		}
		for i := 0; i < len(bySend); i++ {
			bySend[i] = byASCII[iIndex]
		}
		iIndex++
		n, err = uc.WriteToUdpAddr(bySend, addr)
		if err != nil {
			fmt.Printf("write to client[%s] failed: %s", addr.String(), err.Error())
			break
		}
	}
}
