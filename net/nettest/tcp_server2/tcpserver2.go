package main

import (
	"dadian/endian"
	"dadian/net/tcp"
	"fmt"
	"os"
	"time"
)

type E_MSGTYPE int

const (
	MT_PIC  E_MSGTYPE = 1234567890 // 图片消息
	MT_DESC                        // 描述消息
)

type headerT struct {
	iMsgType int32
	iDataLen int32
}
type picInfoT struct {
	iWidth  int32
	iHeight int32
	shortV  int16
	szFile  [128]byte
	ubiSize uint64
}

type picMsgT struct {
	header headerT
	pic    picInfoT
}

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
	err = tsvr.Start(OnConnectS)
	if err != nil {
		fmt.Printf("main: Start failed: %s\n", err.Error())
		return
	}
	for {
		time.Sleep(500 * time.Millisecond)
	}
	tsvr.Stop()
	fmt.Println("Main EXIT")
}

// 获取时间（包括毫秒），返回字符串形式时间：20181203-113059-123
func GetNowString() string {
	var ret string
	tNow := time.Now()
	iNowMS := int(tNow.UnixNano()/1e6) % 1000
	szDate := tNow.Format("2006-01-02 15:04:05")
	ret = fmt.Sprintf("%s-%d", szDate, iNowMS)
	return ret
}

func GetNowString2() string {
	var ret string
	tNow := time.Now()
	iNowMS := int(tNow.UnixNano()/1e6) % 1000
	szDate := tNow.Format("20060102_150405")
	ret = fmt.Sprintf("%s_%03d", szDate, iNowMS)
	return ret
}

func jobS(conn *tcp.TcpConn) {
	// 最大接收缓冲区
	const MAX_BUFF_SIZE = 1024
	byRecvBuf := make([]byte, MAX_BUFF_SIZE)
	//bySendBuf := make([]byte, MAX_BUFF_SIZE)
	var n int = 0
	var err error = nil
	var p []byte = nil
	var intv int
	// 返回时关闭连接
	defer conn.Close()
	szDir := "bmp-pic"
	err = os.Mkdir(szDir, 0777)
	if err != nil {
		fmt.Printf("create dir failed: %s\n", err.Error())
		return
	}
	iIndex := 0
	for {
		iIndex++
		var msg picMsgT
		fmt.Printf("\n###########################################\n%s\n\n", GetNowString())
		// 读取header的消息类型
		p = byRecvBuf[:4]
		n, err = conn.Read(p)
		if err != nil {
			fmt.Printf("job: Read header[msg type] failed: %s\n", err.Error())
			break
		}
		intv, _ = endian.BytesToInt(p, true)
		msg.header.iMsgType = int32(endian.NTOHL(uint32(intv)))
		fmt.Printf("header: msg type = %d\n", msg.header.iMsgType)
		if msg.header.iMsgType != int32(MT_PIC) {
			fmt.Printf("invaild message type: %d", msg.header.iMsgType)
			continue
		}

		// 读取图片信息结构体长度
		p = byRecvBuf[:4]
		n, err = conn.Read(p)
		if err != nil {
			fmt.Printf("job: Read header[data len] failed: %s\n", err.Error())
			break
		}
		intv, _ = endian.BytesToInt(p, true)
		msg.header.iDataLen = int32(endian.NTOHL(uint32(intv)))
		fmt.Printf("heder: pic info len = %d\n", msg.header.iDataLen)

		// 读取图片宽度
		p = byRecvBuf[:4]
		n, err = conn.Read(p)
		if err != nil {
			fmt.Printf("job: Read Picture[width] failed: %s\n", err.Error())
			break
		}
		intv, _ = endian.BytesToInt(p, true)
		msg.pic.iWidth = int32(endian.NTOHL(uint32(intv)))
		fmt.Printf("picture: width = %d\n", msg.pic.iWidth)

		// 读取图片高度
		p = byRecvBuf[:4]
		n, err = conn.Read(p)
		if err != nil {
			fmt.Printf("job: Read Picture[height] failed: %s\n", err.Error())
			break
		}
		intv, _ = endian.BytesToInt(p, true)
		msg.pic.iHeight = int32(endian.NTOHL(uint32(intv)))
		fmt.Printf("picture: height = %d\n", msg.pic.iHeight)

		// 读取图片ID
		p = byRecvBuf[:2]
		n, err = conn.Read(p)
		if err != nil {
			fmt.Printf("job: Read Picture[height] failed: %s\n", err.Error())
			break
		}
		intv, _ = endian.BytesToInt(p, true)
		msg.pic.shortV = int16(endian.NTOHS(uint16(intv)))
		fmt.Printf("picture: shortV = %d\n", msg.pic.shortV)

		// 读取图片文件名称
		p = byRecvBuf[:len(msg.pic.szFile)]
		n, err = conn.Read(p)
		if err != nil {
			fmt.Printf("job: Read Picture[szFile] failed: %s\n", err.Error())
			break
		}
		intv, _ = endian.BytesToInt(p, true)
		copy(msg.pic.szFile[:], p)
		fmt.Printf("picture: szFile = %s\n", string(msg.pic.szFile[:]))

		// 读取图片长度
		p = byRecvBuf[:8]
		n, err = conn.Read(p)
		if err != nil {
			fmt.Printf("job: Read Picture[size] failed: %s\n", err.Error())
			break
		}
		intv, _ = endian.BytesToInt(p, true)
		msg.pic.ubiSize = endian.NTOHLL(uint64(intv))
		fmt.Printf("picture: ubiSize = %d\n", msg.pic.ubiSize)
		ubiPicSize := msg.pic.ubiSize
		if ubiPicSize <= 0 {
			fmt.Printf("picture: size=%d, it's invalid!\n", ubiPicSize)
		}

		// 读取图片数据
		byPicBuf := make([]byte, ubiPicSize)
		n, err = conn.Read(byPicBuf)
		if err != nil {
			fmt.Printf("job: Read Picture data failed: %s\n", err.Error())
		}
		// if err != nil {
		if n != int(ubiPicSize) {
			fmt.Printf("job: Read Picture data failed: read: %d, should be: %d\n", n, ubiPicSize)
		}
		fmt.Printf("picture: read picture data, size = %d\n", n)

		// 写入文件
		szFile := szDir + "/" + fmt.Sprintf("%d_%04d_", msg.pic.shortV, iIndex) + GetNowString2() + ".bmp"

		file, err := os.OpenFile(szFile, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0)
		if err != nil {
			fmt.Printf("job: open|create file[%s] failed: %s\n", szFile, err.Error())
			break
		}
		n, err = file.Write(byPicBuf)
		if err != nil {
			fmt.Printf("write to file[%s] failed: %s\n", szFile, err.Error())
			break
		}
		if n != int(ubiPicSize) {
			fmt.Printf("write to file[%s] ok, but writed length=%d, should be %d\n", szFile, n, ubiPicSize)
		}
		file.Sync()
		file.Close()
		fmt.Printf("write to file[%s] ok, writed size=%d\n", szFile, n)
		fmt.Printf("\n###########################################\n\n")
	} // for {

	fmt.Println("job: finish receive file, exit!")
}
