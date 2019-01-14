package udp

import (
	"testing"
	"time"
)

func TestConnect(t *testing.T) {
	pUdpConn, err := Connect("127.0.0.1", 9001, false)
	if err != nil {
		t.Errorf("%s", err.Error())
		return
	}
	n, err := pUdpConn.Write([]byte("hello, world!"))
	if err != nil {
		t.Errorf("write failed: %s", err.Error())
		return
	}
	byBuf := make([]byte, 1024)
	n, err = pUdpConn.ReadTimeout(byBuf, 5000)
	if err != nil {
		t.Errorf("read failed: %s", err.Error())
		return
	}
	if n == 0 {
		t.Logf("read 0 bytes, exit")
		return
	}
	t.Logf("read from server: size=%d, text=%s\n", n, string(byBuf[:n]))
	time.Sleep(10 * time.Second)
}
