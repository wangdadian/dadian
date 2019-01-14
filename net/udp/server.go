package udp

import (
	"fmt"
	"io"
	"net"
	"time"
)

//
// UDP SERVER
//
//UDP连接信息
type UdpServer struct {
	lastAddr *net.UDPAddr // 最后一次收到数据的远程地址
	conn     *net.UDPConn // 连接句柄
}

// 开启UDP server
// port 监听端口，ip-监听IP，可为空""
// onConnect 当连接发生时回调函数，不可为nil
func NewUdpServer(ip string, port int) (*UdpServer, error) {
	strAddr := ip + fmt.Sprintf(":%d", port)
	udpAddr, err := net.ResolveUDPAddr("udp", strAddr)
	if err != nil {
		return nil, fmt.Errorf("ResolveUDPAddr failed: %s", err.Error())
	}
	pConn, err := net.ListenUDP("udp", udpAddr)
	if err != nil {
		return nil, fmt.Errorf("ListenUDP failed: %s", err.Error())
	}
	pUC := &UdpServer{
		lastAddr: nil,
		conn:     pConn,
	}
	return pUC, nil
}

// 超时读,ms-超时毫秒数
func (self *UdpServer) ReadFromUdpAddrTM(byData []byte, ms int) (int, *net.UDPAddr, error) {
	var err error = nil
	timeout := time.Duration(time.Millisecond.Nanoseconds() * int64(ms))
	if err = self.conn.SetReadDeadline(time.Now().Add(timeout)); err != nil {
		//fmt.Printf("SetReadDeadline failed: %s\n", err.Error())
		return 0, nil, err
	}
	n, addr, err := self.ReadFromUdpAddr(byData)
	self.conn.SetReadDeadline(time.Time{})
	if err != nil {
		return n, addr, err
	}

	return n, addr, nil
}

// 读
func (self *UdpServer) ReadFromUdpAddr(byData []byte) (int, *net.UDPAddr, error) {
	var n int     // 读取的实际长度
	var err error // 读取发生的错误
	var addr *net.UDPAddr
	var addrReal net.UDPAddr
	index := 0
	for index < len(byData) {
		n, addr, err = self.conn.ReadFromUDP(byData[index:])
		// fmt.Printf("\n********** read %d bytes, error: %v, addr: %s\n\n", n, err, addr.String())
		index += n
		if err != nil {
			//
			// 读失败
			//
			// 如果已经没有数据了，不返回失败。
			if err == io.EOF {
				return index, &addrReal, nil
			}
			// 其他非网络错误
			_, ok := err.(net.Error)
			if ok == false {
				return index, &addrReal, err
			} else {
				if addr != nil {
					addrReal = *addr
				}
				return index, &addrReal, nil
			}
			return index, &addrReal, err
		}
		if n == 0 {
			return index, &addrReal, nil
		}
		addrReal = *addr
	}
	return index, &addrReal, nil
}

// 写
func (self *UdpServer) WriteToUdpAddr(byData []byte, addr *net.UDPAddr) (int, error) {
	if len(byData) > 65007 {
		return 0, fmt.Errorf("data is too larger than limit size(65507)")
	}
	var n int     // 写入的实际长度
	var err error // 写入发生的错误
	index := 0
	for index < len(byData) {
		n, err = self.conn.WriteToUDP(byData[index:], addr)
		index += n
		// 写失败
		if err != nil {
			return index, err
		}
	}
	return index, nil
}

// 关闭连接
func (self *UdpServer) Close() error {
	self.conn.Close()
	return nil
}
