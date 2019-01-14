package udp

import (
	"errors"
	"fmt"
	"io"
	"net"
	"strconv"
	"strings"
	"time"
)

// 连接UDP服务器
// serverIP-需要连接的服务器IP地址
// serverPort-需要连接的服务器端口号
// bReConn-是否自动重连
func NewUdpClient(serverIP string, serverPort int, bReConn bool) (*UdpClient, error) {
	szSvrAddr := serverIP + fmt.Sprintf(":%d", serverPort)
	conn, err := net.DialTimeout("udp", szSvrAddr, time.Millisecond*2000)
	if err != nil {
		return nil, fmt.Errorf("Connect to udp server[%s] failed: %s", szSvrAddr, err.Error())
	}
	nu := newUdpClient(conn, bReConn)
	return nu, nil
}

const (
	_                  = iota
	UDP_ERR_NONE       // 网络正常
	UDP_ERR_DISCONNECT // 网络断开
	UDP_ERR_CLOSE      // 网络关闭
)

//UDP连接信息
type UdpClient struct {
	LIP         string          // 本地IP
	LPort       int             // 本地连接端口
	RIP         string          // 远程IP
	RPort       int             // 远程连接端口
	conn        net.Conn        // 连接句柄
	bReconn     bool            // 是否重连
	cbException func(err error) // 异常回调
	connState   int             // 网络连接状态
	bExit       bool            // 退出通道
}

func newUdpClient(conn net.Conn, bReConn bool) *UdpClient {
	if conn == nil {
		return nil
	}

	szLAddr := conn.LocalAddr().String()
	szLIP := szLAddr[:strings.Index(szLAddr, ":")]
	szLPort := szLAddr[strings.Index(szLAddr, ":")+1:]
	iLPort, _ := strconv.Atoi(szLPort)
	szRIP := ""
	szRPort := ""
	iRPort := 0
	if conn.RemoteAddr() != nil {
		szRAddr := conn.RemoteAddr().String()
		szRIP = szRAddr[:strings.Index(szRAddr, ":")]
		szRPort = szRAddr[strings.Index(szRAddr, ":")+1:]
		iRPort, _ = strconv.Atoi(szRPort)
	}

	var pUC *UdpClient = &UdpClient{
		conn:        conn,
		LIP:         szLIP,
		LPort:       iLPort,
		RIP:         szRIP,
		RPort:       iRPort,
		cbException: nil,
		connState:   UDP_ERR_NONE,
		bReconn:     bReConn,
		bExit:       false,
	}
	if pUC.bReconn {
		go pUC.thReConnect()
	}
	return pUC
}

// 定时检测连接正常与否
func (self *UdpClient) thReConnect() {
	//fmt.Println("thReConnect: start!")
	for {
		if self.bExit {
			break
		}

		// 网络正常
		if self.connState == UDP_ERR_NONE {
			//fmt.Println("it's ok, continue thIsOk!")
			time.Sleep(500 * time.Millisecond)
			continue
		}
		// 网络异常
		// 清理数据
		self.conn.Close()
		self.conn = nil
		// 客户端连接类型，重连机制
		if self.bReconn == true {
			self.reConnectForClient()
			time.Sleep(500 * time.Millisecond)
			continue
		}
		break
	}
}

// 客户端连接模式下的重连
func (self *UdpClient) reConnectForClient() error {
	conn, err := net.DialTimeout("udp", fmt.Sprintf("%s:%d", self.RIP, self.RPort), 2000*time.Millisecond)
	if err != nil {
		return err
	}
	self.conn = conn
	szLAddr := conn.LocalAddr().String()
	szLIP := szLAddr[:strings.Index(szLAddr, ":")]
	szLPort := szLAddr[strings.Index(szLAddr, ":")+1:]
	iLPort, _ := strconv.Atoi(szLPort)
	self.LIP = szLIP
	self.LPort = iLPort
	// 修改连接状态为正常
	self.connState = UDP_ERR_NONE
	return nil
}

// 设置连接异常时回调函数
func (self *UdpClient) SetOnExceptionCB(cbException func(err error)) bool {
	self.cbException = cbException
	return true
}

// 判断网络是否正常
func (self *UdpClient) IsOK() bool {
	return self.connState == UDP_ERR_NONE
}

// 超时读,ms-超时毫秒数
func (self *UdpClient) ReadTimeout(byData []byte, ms int) (int, error) {
	var err error = nil
	timeout := time.Duration(time.Millisecond.Nanoseconds() * int64(ms))
	if err = self.conn.SetReadDeadline(time.Now().Add(timeout)); err != nil {
		//fmt.Printf("SetReadDeadline failed: %s\n", err.Error())
		return 0, err
	}
	n, err := self.Read(byData)
	self.conn.SetReadDeadline(time.Time{})
	if err != nil {
		return n, err
	}

	return n, nil
}

// 读
func (self *UdpClient) Read(byData []byte) (int, error) {
	if !self.IsOK() {
		return -1, errors.New("net is not ok")
	}

	var n int     // 读取的实际长度
	var err error // 读取发生的错误
	index := 0
	for index < len(byData) {
		n, err = self.conn.Read(byData[index:])
		index += n
		if err != nil {
			//
			// 读失败
			//
			// 如果已经没有数据了，不返回失败。
			if err == io.EOF {
				return index, nil
			}

			// 其他非网络错误
			_, ok := err.(net.Error)
			if ok == false {
				self.connState = UDP_ERR_DISCONNECT
				// 异常回调
				if self.cbException != nil {
					self.cbException(err)
				}
				return index, err
			} else {
				return index, nil
			}
		}
		if n == 0 {
			return index, nil
		}
	}
	return index, nil
}

// 写
func (self *UdpClient) Write(byData []byte) (int, error) {
	if !self.IsOK() {
		return -1, errors.New("net is not ok")
	}
	var n int     // 写入的实际长度
	var err error // 写入发生的错误
	index := 0
	for index < len(byData) {
		n, err = self.conn.Write(byData[index:])
		index += n
		// 写失败
		if err != nil {
			self.connState = UDP_ERR_DISCONNECT
			// 异常回调
			if self.cbException != nil {
				self.cbException(err)
			}
			return index, err
		}
	}
	return index, nil
}

// 关闭连接
func (self *UdpClient) Close() error {
	self.close()
	return nil
}

func (self *UdpClient) close() {
	// 停止检测连接状态的 goroutine
	self.bExit = true
	// 关闭连接
	if self.IsOK() {
		self.conn.Close()
	}

	// 清理数据
	self.LIP = ""
	self.LPort = 0
	self.RIP = ""
	self.RPort = 0
	// 修改连接状态为断开
	self.connState = UDP_ERR_DISCONNECT
	// 关闭连接后禁用自动重连
	self.bReconn = false
}
