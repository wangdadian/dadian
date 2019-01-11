package tcp

import (
	"errors"
	// "fmt"
	"io"
	"net"
	"strconv"
	"strings"
	"time"
)

const (
	_                  = iota
	TCP_ERR_NONE       // 网络连接正常
	TCP_ERR_DISCONNECT // 网络连接断开
	TCP_ERR_CLOSE      // 网络连接关闭
)

//TCP连接信息
type TcpConn struct {
	LIP         string          // 本地IP
	LPort       int             // 本地连接端口
	RIP         string          // 远程IP
	RPort       int             // 远程连接端口
	conn        net.Conn        // 连接句柄
	bSvrTcp     bool            // tcp server端连接
	bReconn     bool            // 是否重连
	cbExpection func(err error) // 异常回调
	connState   int             // 网络连接状态
	exit        chan bool       // 退出通道
	data        interface{}     // 用户数据

}

// 新建tcp连接信息
func newTcpConn(conn net.Conn, bSvrTcp bool, bReConn bool, data interface{}) *TcpConn {
	if conn == nil {
		return nil
	}
	szLAddr := conn.LocalAddr().String()
	szRAddr := conn.RemoteAddr().String()
	szLIP := szLAddr[:strings.Index(szLAddr, ":")]
	szLPort := szRAddr[strings.Index(szRAddr, ":")+1:]
	iLPort, _ := strconv.Atoi(szLPort)
	szRIP := szRAddr[:strings.Index(szRAddr, ":")]
	szRPort := szRAddr[strings.Index(szRAddr, ":")+1:]
	iRPort, _ := strconv.Atoi(szRPort)
	var tc *TcpConn = &TcpConn{
		conn:        conn,
		LIP:         szLIP,
		LPort:       iLPort,
		RIP:         szRIP,
		RPort:       iRPort,
		cbExpection: nil,
		connState:   TCP_ERR_NONE,
		bReconn:     bReConn,
		bSvrTcp:     bSvrTcp,
		data:        data,
		exit:        make(chan bool),
	}
	go tc.thIsOK()
	return tc
}

// 定时检测连接正常与否
func (self *TcpConn) thIsOK() {
	//fmt.Println("thIsOK: start!")
	for {
		select {
		case v := <-self.exit:
			if v {
				goto goto_exit
			}
			break
		default:
			//fmt.Println("it's select default!")
		}

		// 网络正常
		if self.connState == TCP_ERR_NONE {
			//fmt.Println("it's ok, continue thIsOk!")
			time.Sleep(500 * time.Millisecond)
			continue
		}
		// 网络异常
		// 清理数据
		self.conn.Close()
		self.conn = nil
		// 客户端连接类型，重连机制
		if self.bSvrTcp == false && self.bReconn == true {
			// 类型断言
			if v, ok := self.data.(TcpClientInfo); ok {
				var tcpClient TcpClientInfo = v
				reConnect_Client(&tcpClient, self)
				time.Sleep(500 * time.Millisecond)
				continue
			}
		}
		break
	}
goto_exit:
	//fmt.Println("thIsOK: exit!")
}

// 客户端连接模式下的重连
func reConnect_Client(tc *TcpClientInfo, tcpconn *TcpConn) error {
	conn, err := net.DialTimeout("tcp", tc.addr, 2000*time.Millisecond)
	if err != nil {
		return err
	}
	tcpconn.conn = conn
	szLAddr := conn.LocalAddr().String()
	szRAddr := conn.RemoteAddr().String()
	szLIP := szLAddr[:strings.Index(szLAddr, ":")]
	szLPort := szRAddr[strings.Index(szRAddr, ":")+1:]
	iLPort, _ := strconv.Atoi(szLPort)
	szRIP := szRAddr[:strings.Index(szRAddr, ":")]
	szRPort := szRAddr[strings.Index(szRAddr, ":")+1:]
	iRPort, _ := strconv.Atoi(szRPort)
	tcpconn.LIP = szLIP
	tcpconn.LPort = iLPort
	tcpconn.RIP = szRIP
	tcpconn.RPort = iRPort
	// 修改连接状态为正常
	tcpconn.connState = TCP_ERR_NONE
	return nil
}

// 设置连接异常时回调函数
func (self *TcpConn) SetOnExpectionCB(cbExpection func(err error)) bool {
	self.cbExpection = cbExpection
	return true
}

// 判断网络是否正常
func (self *TcpConn) IsOK() bool {
	return self.connState == TCP_ERR_NONE
}

// 超时读,ms-超时毫秒数
func (self *TcpConn) ReadTimeout(byData []byte, ms int) (int, error) {
	var err error = nil
	timeout := time.Duration(time.Millisecond.Nanoseconds() * int64(ms))
	if err = self.conn.SetReadDeadline(time.Now().Add(timeout)); err != nil {
		//fmt.Printf("SetReadDeadline failed: %s\n", err.Error())
		return 0, nil
	}
	n, err := self.Read(byData)
	self.conn.SetReadDeadline(time.Time{})
	if err != nil {
		return n, err
	}

	return n, nil
}

// 读
func (self *TcpConn) Read(byData []byte) (int, error) {
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
				self.connState = TCP_ERR_DISCONNECT
				// 异常回调
				if self.cbExpection != nil {
					self.cbExpection(err)
				}
				return index, err
			} else {
				return index, nil
			}
		}
	}
	return index, nil
}

// 写
func (self *TcpConn) Write(byData []byte) (int, error) {
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
			self.connState = TCP_ERR_DISCONNECT
			// 异常回调
			if self.cbExpection != nil {
				self.cbExpection(err)
			}
			return index, err
		}
	}
	return index, nil
}

// 关闭连接
func (self *TcpConn) Close() error {
	self.close()
	return nil
}

func (self *TcpConn) close() {
	// 停止检测连接状态的 goroutine
	self.exit <- true
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
	self.connState = TCP_ERR_DISCONNECT
	// 关闭连接后禁用自动重连
	self.bReconn = false
}
