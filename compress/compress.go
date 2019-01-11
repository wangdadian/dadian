package compress

import (
	"dadian/compress/targz"
	"dadian/compress/zip"
	"errors"
	"fmt"
	"path"
	"runtime"
	"strings"
)

/*
 * 已知bug：
 *   1、空目录压缩后会丢失， 2018-01-04
 *
 */

type CompressType int

const (
	_        CompressType = iota // 无效
	CT_AUTO                      // 根据os类型选择压缩模式，支持linux、windows
	CT_TARGZ                     // "tar.gz" 模式，linux默认模式
	CT_ZIP                       // ".zip"模式，windows默认模式，以及其他系统默认模式
)

type CompressIntf interface {

	// 压缩文件
	Compress(szFiles []string, szDestFile string) (string, error)

	// 解压文件
	DeCompress(srcFile, destDir string) error
}

// 压缩文件
// szFiles:需要压缩的文件列表
// szDestFile：压缩后的目标文件名称，会根据压缩模式自动添加后缀名称
// type 压缩模式
// 返回： 成功返回目标文件名称、nil，失败返回""、错误信息
func Compress(szFiles []string, szDestFile string, ztype CompressType) (string, error) {
	var pCompressIntf CompressIntf
	switch ztype {
	case CT_AUTO:
		if runtime.GOOS == "linux" {
			pCompressIntf = &targz.TargzCompress{}
		} else if runtime.GOOS == "windows" {
			pCompressIntf = &zip.ZipCompress{}
		} else {
			pCompressIntf = &zip.ZipCompress{}
		}
	case CT_TARGZ:
		pCompressIntf = &targz.TargzCompress{}
	case CT_ZIP:
		pCompressIntf = &zip.ZipCompress{}
	default:
		return "", errors.New(fmt.Sprintf("unsupported compress type[%d]", ztype))
	}
	return pCompressIntf.Compress(szFiles, szDestFile)
}

// 解压文件
// 根据后缀名自动判断解压模式
// srcFile: 待解压文件
// destDir：解压目标路径
// 返回：成功返回nil，失败返回错误信息
func DeCompress(srcFile, destDir string) error {
	var pCompressIntf CompressIntf
	szExt := path.Ext(srcFile)
	szExt = strings.ToLower(szExt)
	switch szExt {
	case ".zip":
		pCompressIntf = &zip.ZipCompress{}
	case ".gz", ".tar.gz":
		pCompressIntf = &targz.TargzCompress{}
	default:
		return errors.New(fmt.Sprintf("unsupported to decompress file[%s]", srcFile))
	}
	return pCompressIntf.DeCompress(srcFile, destDir)
}
