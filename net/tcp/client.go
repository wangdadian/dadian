package tcp

import (
	"fmt"
	"net"
	"time"
)

// 连接服务器
// bReconn 连接成功后的时间里，如果出现断开是否重连
// onConnect 当连接发生时回调函数，不可为nil
func Connect(ip string, port int, bReconn bool) (*TcpConn, error) {
	conn, err := net.DialTimeout("tcp", fmt.Sprintf("%s:%d", ip, port), 2000*time.Millisecond)
	if err != nil {
		return nil, fmt.Errorf("Dial [%s:%d] failed: %s", ip, port, err.Error())
	}
	// 连接成功，创建连接操作句柄，并回调传给上层
	tc := newTcpConn(conn, bReconn)
	return tc, nil
}
