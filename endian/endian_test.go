package endian

import (
	"runtime"
	"testing"
)

func TestEdian1(t *testing.T) {
	// 本机大小端信息
	edian := GetEndian()
	switch edian {
	case ENDIAN_BIG:
		t.Logf("SYS %s: Big Endian!", runtime.GOOS)
		break
	case ENDIAN_LITTLE:
		t.Logf("SYS %s: Little Endian!", runtime.GOOS)
		break
	default:
		t.Errorf("SYS %s: unkown endian!", runtime.GOOS)
	}

	var sh uint16 = 12345
	shN := HTONS(sh)
	t.Logf("HTONS: %d [0x%04x] -> %d[0x%04x]\n", sh, sh, shN, shN)
	shH := NTOHS(sh)
	t.Logf("NTOHS: %d [0x%04x] -> %d[0x%04x]\n", sh, sh, shH, shH)

	var ui uint32 = 12345
	uiN := HTONL(ui)
	t.Logf("HTONL: %d [0x%04x] -> %d[0x%04x]\n", ui, ui, uiN, uiN)
	uiH := NTOHL(ui)
	t.Logf("NTOHL: %d [0x%04x] -> %d[0x%04x]\n", ui, ui, uiH, uiH)

	var ubi uint64 = 12345
	ubiN := HTONLL(ubi)
	t.Logf("HTONLL: %d [0x%04x] -> %d[0x%04x]\n", ubi, ubi, ubiN, ubiN)
	ubiH := NTOHLL(ubi)
	t.Logf("NTOHLL: %d [0x%04x] -> %d[0x%04x]\n", ubi, ubi, ubiH, ubiH)
}

func TestByteToInt(t *testing.T) {
	var b []byte = []byte{0x11, 0x22, 0x33, 0x44, 0x55, 0x66, 0x67, 0x68}
	n32, err := BytesToInt(b[:4], true)
	if err != nil {
		t.Errorf("BytesToInt: %s", err.Error())
	}
	t.Logf("BytesToInt: b[:4] = %d[%x]", n32, n32)

	n16, err := BytesToInt(b[5:7], false)
	if err != nil {
		t.Errorf("BytesToInt: %s", err.Error())
	}
	t.Logf("BytesToInt: b[5:7] = %d[%x]", n16, n16)

	n64, err := BytesToInt(b[:], false)
	if err != nil {
		t.Errorf("BytesToInt: %s", err.Error())
	}
	t.Logf("BytesToInt: b[:] = %d[%x]", n64, n64)
}

func TestIntToBytes(t *testing.T) {
	edian := GetEndian()
	switch edian {
	case ENDIAN_BIG:
		t.Logf("ENDIAN_BIG")
	case ENDIAN_LITTLE:
		t.Logf("ENDIAN_LITTLE")
	default:
		t.Errorf("ENDIAN_LITTLE")
	}

	var n int64 = 0x1877665544332211
	b, err := IntToBytes(int(n), 8)
	if err != nil {
		t.Errorf("IntToBytes failed: %s", err.Error())
	}
	for i, v := range b {
		t.Logf("%d:%x  ", i, v)
	}
}
