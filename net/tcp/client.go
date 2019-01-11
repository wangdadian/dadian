package tcp

import (
	"errors"
	"fmt"
	"net"
	"time"
)

type TcpClientInfo struct {
	Port int    // 需要连接的服务器端口
	Ip   string // 需要连接的服务器IP
	addr string // 服务器IP:端口
}

func NewTcpClient(ip string, port int) (*TcpClientInfo, error) {
	byIP := net.ParseIP(ip)
	if port <= 0 || len(ip) < 7 || byIP == nil {
		return nil, errors.New(fmt.Sprintf("invaild port[%d] or ip[%s]", port, ip))
	}
	c := &TcpClientInfo{
		Port: port,
		Ip:   ip,
		addr: fmt.Sprintf("%s:%d", ip, port),
	}
	return c, nil
}

// 连接服务器
// bReconn 连接成功后的时间里，如果出现断开是否重连
// onConnect 当连接发生时回调函数，不可为nil
func (self *TcpClientInfo) Connect(bReconn bool, onConnect func(*TcpConn)) error {
	if onConnect == nil {
		return errors.New("onConnect handler is nil.")
	}
	conn, err := net.DialTimeout("tcp", self.addr, 2000*time.Millisecond)
	if err != nil {
		return err
	}
	// 连接成功，创建连接操作句柄，并回调传给上层
	tc := newTcpConn(conn, false, bReconn, self)
	go onConnect(tc)
	return nil
}
