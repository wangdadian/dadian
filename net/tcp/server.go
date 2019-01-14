package tcp

import (
	"errors"
	"fmt"
	"net"
	"time"
	// "sync"
)

type TcpSvrInfo struct {
	Port int
	Ip   string
	addr string
	lsn  net.Listener
}

//
// TCP SERVER
//

// 新建TCPserver对象
// port 监听端口，ip-监听IP，可为空""
func NewTcpServer(ip string, port int) (*TcpSvrInfo, error) {
	if port < 0 {
		return nil, errors.New("invalid port")
	}
	tcpSvr := TcpSvrInfo{
		Port: port,
		Ip:   ip,
		addr: fmt.Sprintf("%s:%d", ip, port),
		lsn:  nil,
	}

	return &tcpSvr, nil
}

// 开启TCP server
// onConnect 当连接发生时回调函数，不可为nil
func (self *TcpSvrInfo) Start(onConnect func(*TcpConn)) error {
	if onConnect == nil {
		return errors.New("invalid hanldle of onConnect")
	}
	listenMe, err := net.Listen("tcp", self.addr)
	if err != nil {
		return err
	}
	//fmt.Printf("Listen at port %d, ok.\n", self.Port)
	self.lsn = listenMe
	// 给通道分配空间，防止Accept成功后阻塞导致通道阻塞
	errch := make(chan error, 1)
	defer close(errch)

	// 开始等待连接
	go self.accept(errch, onConnect)

	// 睡眠
	time.Sleep(time.Millisecond * 10)
	errAccept := <-errch

	if errAccept != nil {
		return errAccept
	}
	return nil
}

func (self *TcpSvrInfo) accept(errch chan error, onConnect func(*TcpConn)) error {
	bFirst := true
	// 赋予通道一个nil值，保证上层读取时有值可取，避免阻塞
	errch <- nil
	for {
		//fmt.Printf("Waiting for connection.\n")
		conn, err := self.lsn.Accept()
		if err != nil {
			if bFirst {
				// 取出通道中的nil值，赋予错误值，首先判断chan是否已被close
				if _, ok := <-errch; !ok {
					break
				}
				errch <- err
				bFirst = false
			}
			//fmt.Printf("server[%s], Accept failed: %s", self.addr, err.Error())
			continue
		}

		//fmt.Printf("#### server[%s], new client: %s\n", self.addr, conn.RemoteAddr().String())
		tc := newTcpConn(conn, false)
		// 连接时回调
		go onConnect(tc)
	}
	return nil
}

// 服务停止
// onStop - 停止服务时的回调函数，可谓nil
func (self *TcpSvrInfo) Stop(onStop func()) error {
	if onStop != nil {
		go onStop()
	}
	err := self.lsn.Close()
	if err != nil {
		return err
	}
	self.lsn = nil
	return nil
}
