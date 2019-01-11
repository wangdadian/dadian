package main

import (
	"dadian/compress"
	"fmt"
	"io/ioutil"
	"runtime"
	"time"
)

func main() {
	// 待压缩的文件夹
	path := "/root/test/"
	if runtime.GOOS == "windows" {
		path = "D:\\golang\\src\\test\\"
	}
	// 压缩到目标文件，自动添加后缀名
	szDestFile := "/tmp/testz"
	if runtime.GOOS == "windows" {
		szDestFile = "d:\\temp\\testz"
	}
	// 解压目标文件夹
	szDir := "/tmp/test"
	if runtime.GOOS == "windows" {
		szDir = "d:\\temp\\test"
	}

	// 获取待压缩的文件夹下所有子文件夹及文件信息
	szFiles := []string{}
	files, _ := ioutil.ReadDir(path)
	for _, f := range files {
		szFiles = append(szFiles, path+f.Name())
	}

	// 压缩
	dest, err := compress.Compress(szFiles, szDestFile, compress.CT_TARGZ)
	if err != nil {
		fmt.Printf("compress failed: %s\n", err.Error())
		return
	}
	fmt.Printf("compress ok, target file: %s\n", dest)

	// 解压
	if err := compress.DeCompress(dest, szDir); err != nil {
		fmt.Printf("decompress [%s] failed: %s\n", dest, err.Error())
	} else {
		fmt.Printf("decompress file [%s] to %s ok\n", dest, szDir)
	}

	// sleep
	for {
		time.Sleep(1 * time.Second)
	}
}
