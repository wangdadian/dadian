/*
 * @Author: wangdadian
 * @Date: 2018-12-27 14:26:10
 * @Last Modified by: wangdadian
 * @Last Modified time: 2018-12-27 15:00:02
 */

/*
 * 1、获取当前机器大小端类别
 * 2、网络、本机字节序转换
 * 3、字节流转换成本机数值（依据本机大小端类型）
 * 4、数值转换成字节流（依据本机大小端类型）
 */

package endian

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"unsafe"
)

// 机器大小端
type SYSENDIAN int

// 机器大小端类型
const (
	_             SYSENDIAN = iota
	ENDIAN_BIG              // 大端
	ENDIAN_LITTLE           // 小端
)

// 获取当前机器大小端类型
func GetEndian() SYSENDIAN {
	const N int = int(unsafe.Sizeof(0))
	x := 0x00001234
	p := unsafe.Pointer(&x)
	bs := (*[N]byte)(p)
	if bs[0] == 0 {

		return ENDIAN_BIG
	} else {
		return ENDIAN_LITTLE
	}
}

// 网络大端转换成本机数值，16位整型
func NTOHS(n uint16) uint16 {
	retn := n
	edian := GetEndian()
	// 本机为大端，则直接继承网络的大端数据值
	if edian == ENDIAN_BIG {
		return retn
	}
	// 本机为小端，则转换成本机的小端模式
	b := make([]byte, 2)
	// 原数值为大端值，取出大端字节流
	binary.BigEndian.PutUint16(b, n)
	// 按照小端模式读取度去除大端字节流
	bb := bytes.NewBuffer(b)
	binary.Read(bb, binary.LittleEndian, &retn)
	return retn
}

// 本机数值转换成网络大端，16位整型
func HTONS(n uint16) uint16 {
	retn := n
	edian := GetEndian()
	// 本机为大端，则直接赋值给网络的大端数据
	if edian == ENDIAN_BIG {
		return retn
	}

	// 本机为小端，则转换成网络的大端模式
	b := make([]byte, 2)
	// 原数值为小端数据，取出小端字节流
	binary.LittleEndian.PutUint16(b, n)
	bb := bytes.NewBuffer(b)
	// 按照大端模式读取出小端字节流
	binary.Read(bb, binary.BigEndian, &retn)
	return retn
}

// 网络大端转换成本机数值，32位整型
func NTOHL(n uint32) uint32 {
	retn := n
	edian := GetEndian()
	// 本机为大端，则直接赋值给网络的大端数据
	if edian == ENDIAN_BIG {
		return retn
	}

	// 本机为小端，则转换成小端模式
	b := make([]byte, 4)
	// 原数值为大端数据，取出大端字节流
	binary.BigEndian.PutUint32(b, n)
	bb := bytes.NewBuffer(b)
	// 按照小端端模式读取出大端字节流
	binary.Read(bb, binary.LittleEndian, &retn)
	return retn
}

// 本机数值转换成网络大端，32位整型
func HTONL(n uint32) uint32 {
	retn := n
	edian := GetEndian()
	// 本机为大端，则直接赋值给网络的大端数据
	if edian == ENDIAN_BIG {
		return retn
	}

	// 本机为小端，则转换成网络大端模式
	b := make([]byte, 4)
	// 原数值为小端数据，取出小端字节流
	binary.LittleEndian.PutUint32(b, n)
	bb := bytes.NewBuffer(b)
	// 按照大端端端模式读取出小端字节流
	binary.Read(bb, binary.BigEndian, &retn)
	return retn
}

// 网络大端转换成本机数值，64位整型
func NTOHLL(n uint64) uint64 {
	retn := n
	edian := GetEndian()
	// 本机为大端，则直接赋值给网络的大端数据
	if edian == ENDIAN_BIG {
		return retn
	}

	// 本机为小端，则转换成小端模式
	b := make([]byte, 8)
	// 原数值为大端数据，取出大端字节流
	binary.BigEndian.PutUint64(b, n)
	bb := bytes.NewBuffer(b)
	// 按照小端端模式读取出大端字节流
	binary.Read(bb, binary.LittleEndian, &retn)
	return retn
}

// 本机数值转换成网络大端，64位整型
func HTONLL(n uint64) uint64 {
	retn := n
	edian := GetEndian()
	// 本机为大端，则直接赋值给网络的大端数据
	if edian == ENDIAN_BIG {
		return retn
	}

	// 本机为小端，则转换成网络大端模式
	b := make([]byte, 8)
	// 原数值为小端数据，取出小端字节流
	binary.LittleEndian.PutUint64(b, n)
	bb := bytes.NewBuffer(b)
	// 按照大端端端模式读取出小端字节流
	binary.Read(bb, binary.BigEndian, &retn)
	return retn
}

//整形转换成字节
// n-待转换的整型
// b-转换成的字节数，如int32:4, int8:1
func IntToBytes(n int, b byte) ([]byte, error) {
	bBigEndian := true
	edian := GetEndian()
	switch edian {
	case ENDIAN_BIG:
		bBigEndian = true
		break
	case ENDIAN_LITTLE:
		bBigEndian = false
		break
	default:
	}

	switch b {
	case 1:
		tmp := int8(n)
		bytesBuffer := bytes.NewBuffer([]byte{})
		if bBigEndian {
			binary.Write(bytesBuffer, binary.BigEndian, &tmp)
		} else {
			binary.Write(bytesBuffer, binary.LittleEndian, &tmp)
		}
		return bytesBuffer.Bytes(), nil
	case 2:
		tmp := int16(n)
		bytesBuffer := bytes.NewBuffer([]byte{})
		if bBigEndian {
			binary.Write(bytesBuffer, binary.BigEndian, &tmp)
		} else {
			binary.Write(bytesBuffer, binary.LittleEndian, &tmp)
		}
		return bytesBuffer.Bytes(), nil
	case 4:
		tmp := int32(n)
		bytesBuffer := bytes.NewBuffer([]byte{})
		if bBigEndian {
			binary.Write(bytesBuffer, binary.BigEndian, &tmp)
		} else {
			binary.Write(bytesBuffer, binary.LittleEndian, &tmp)
		}
		return bytesBuffer.Bytes(), nil
	case 8:
		tmp := int64(n)
		bytesBuffer := bytes.NewBuffer([]byte{})
		if bBigEndian {
			binary.Write(bytesBuffer, binary.BigEndian, &tmp)
		} else {
			binary.Write(bytesBuffer, binary.LittleEndian, &tmp)
		}
		return bytesBuffer.Bytes(), nil
	}
	return nil, fmt.Errorf("IntToBytes b param is invaild")
}

//
// 字节流转换成本机数值
// b	字节流数据，长度必须为1,2,4,8，其他不支持并返回错误，
//      根据 b的长度匹配相应类型，如长度为2表示为16位整型，4-32位整型，8-64位整型
// isSymbol	表示有无符号
// 统一返回int型数值，外部调用进行类型转换
func BytesToInt(b []byte, isSymbol bool) (int, error) {
	bBigEndian := true
	edian := GetEndian()
	switch edian {
	case ENDIAN_BIG:
		bBigEndian = true
		break
	case ENDIAN_LITTLE:
		bBigEndian = false
		break
	default:
	}

	if isSymbol {
		return bytesToIntS(b, bBigEndian)
	}
	return bytesToIntU(b, bBigEndian)
}

//字节数组转成int(无符号的)
func bytesToIntU(b []byte, bBigEndian bool) (int, error) {
	bytesBuffer := bytes.NewBuffer(b)
	var err error = nil
	switch len(b) {
	case 1:
		var tmp uint8
		if bBigEndian {
			err = binary.Read(bytesBuffer, binary.BigEndian, &tmp)
		} else {
			err = binary.Read(bytesBuffer, binary.LittleEndian, &tmp)
		}
		return int(tmp), err
	case 2:
		var tmp uint16
		if bBigEndian {
			err = binary.Read(bytesBuffer, binary.BigEndian, &tmp)
		} else {
			err = binary.Read(bytesBuffer, binary.LittleEndian, &tmp)
		}
		return int(tmp), err
	case 4:
		var tmp uint32
		if bBigEndian {
			err = binary.Read(bytesBuffer, binary.BigEndian, &tmp)
		} else {
			err = binary.Read(bytesBuffer, binary.LittleEndian, &tmp)
		}
		return int(tmp), err
	case 8:
		var tmp uint64
		if bBigEndian {
			err = binary.Read(bytesBuffer, binary.BigEndian, &tmp)
		} else {
			err = binary.Read(bytesBuffer, binary.LittleEndian, &tmp)
		}
		return int(tmp), err
	default:
		return 0, fmt.Errorf("%s", "BytesToInt: bytes length is invaild.")
	}
}

//字节数组转成int(有符号)
func bytesToIntS(b []byte, bBigEndian bool) (int, error) {
	bytesBuffer := bytes.NewBuffer(b)
	var err error = nil
	switch len(b) {
	case 1:
		var tmp int8
		if bBigEndian {
			err = binary.Read(bytesBuffer, binary.BigEndian, &tmp)
		} else {
			err = binary.Read(bytesBuffer, binary.LittleEndian, &tmp)
		}
		return int(tmp), err
	case 2:
		var tmp int16
		if bBigEndian {
			err = binary.Read(bytesBuffer, binary.BigEndian, &tmp)
		} else {
			err = binary.Read(bytesBuffer, binary.LittleEndian, &tmp)
		}
		return int(tmp), err
	case 4:
		var tmp int32
		if bBigEndian {
			err = binary.Read(bytesBuffer, binary.BigEndian, &tmp)
		} else {
			err = binary.Read(bytesBuffer, binary.LittleEndian, &tmp)
		}
		return int(tmp), err
	case 8:
		var tmp int64
		if bBigEndian {
			err = binary.Read(bytesBuffer, binary.BigEndian, &tmp)
		} else {
			err = binary.Read(bytesBuffer, binary.LittleEndian, &tmp)
		}
		return int(tmp), err
	default:
		return 0, fmt.Errorf("%s", "BytesToInt: bytes length is invaild.")
	}
}
