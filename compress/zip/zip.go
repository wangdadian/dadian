package zip

import (
	"archive/zip"
	"errors"
	"fmt"
	"io"
	"os"
	"runtime"
	"strings"
)

type ZipCompress struct {
}

func (self *ZipCompress) Compress(szFiles []string, szDestFile string) (string, error) {
	// 目标文件非目录
	if strings.HasSuffix(szDestFile, "/") || strings.HasSuffix(szDestFile, "\\") {
		return "", errors.New("dest file name is a dir")
	}
	if yes := strings.HasSuffix(szDestFile, ".zip"); !yes {
		szDestFile += ".zip"
	}
	// 创建目标文件
	df, err := os.Create(szDestFile)
	if err != nil {
		return "", errors.New(fmt.Sprintf("create dest file[%s] failed: %s", szDestFile, err.Error()))
	}
	defer df.Close()
	zw := zip.NewWriter(df)
	defer zw.Close()
	// 打开文件，并压缩
	for _, szFile := range szFiles {
		f, err := os.Open(szFile)
		if err != nil {
			err = errors.New(fmt.Sprintf("open file[%s] failed: %s", szFile, err.Error()))
			goto goto_ret
		}
		if err := compress(f, "", zw); err != nil {
			goto goto_ret
		}
	}
goto_ret:
	// 压缩失败，删除已创建的压缩文件
	if err != nil {
		os.Remove(szDestFile)
		return "", err
	}
	return szDestFile, nil
}

// 解压文件
func (self *ZipCompress) DeCompress(zipFile, destDir string) error {
	srcFile, err := os.Open(zipFile)
	if err != nil {
		return fmt.Errorf("os.Open file[%s] failed: %s", zipFile, err.Error())
	}
	defer srcFile.Close()
	zr, err := zip.OpenReader(zipFile)
	if err != nil {
		return fmt.Errorf("zip.OpenReader file[%s] failed: %s", zipFile, err.Error())
	}
	defer zr.Close()
	for _, f := range zr.File {
		path := destDir
		if runtime.GOOS == "linux" {
			if yes := strings.HasSuffix(destDir, "/"); !yes {
				path += "/"
			}
		} else if runtime.GOOS == "windows" {
			if yes := strings.HasSuffix(destDir, "\\"); !yes {
				path += "\\"
			}
		} else {
			if yes := strings.HasSuffix(destDir, "/"); !yes {
				path += "/"
			}
		}
		filename := path + f.Name
		info := f.FileInfo()
		// 如果是文件夹则创建，并继续处理下一个文件（夹）信息
		if info.IsDir() {
			err = os.MkdirAll(filename, 0755)
			if err != nil {
				return fmt.Errorf("create dir [%s] failed: %s", filename, err.Error())
			}
			continue
		}
		// 文件信息处理
		zipF, err := f.Open()
		if err != nil {
			return fmt.Errorf("zip.File.Open [%s] failed: %s", f.Name, err.Error())
		}
		defer zipF.Close()
		unzipF, err := createFile(filename)
		if err != nil {
			return err
		}
		_, err = io.Copy(unzipF, zipF)
		unzipF.Close()
		if err != nil {
			return fmt.Errorf("io.Copy from [%s] to [%s] failed: %s", f.Name, unzipF.Name(), err.Error())
		}
	}

	return nil
}

func createFile(name string) (*os.File, error) {
	var err error = nil
	var dir string
	if runtime.GOOS == "linux" {
		dir = string([]rune(name)[0:strings.LastIndex(name, "/")])
	} else if runtime.GOOS == "windows" {
		dir = string([]rune(name)[0:strings.LastIndex(name, "\\")])
	} else {
		dir = string([]rune(name)[0:strings.LastIndex(name, "/")])
	}
	err = os.MkdirAll(dir, 0755)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("create dir [%s] failed: %s", dir, err.Error()))
	}
	return os.Create(name)
}

// 压缩文件、递归压缩文件夹及其子目录中的所有文件，自动关闭传入的file
func compress(file *os.File, prefix string, zw *zip.Writer) error {
	// fmt.Printf("compress: file=%s\n", file.Name())
	defer func() {
		err := file.Close()
		if err != nil {
			fmt.Printf("close file[%s] failed: %s\n", file.Name(), err.Error())
		}
	}()
	info, err := file.Stat()
	if err != nil {
		return fmt.Errorf("os.File.Stat failed: %s", err.Error())
	}

	if info.IsDir() {
		if prefix == "" {
			prefix = info.Name()
		} else {
			if runtime.GOOS == "linux" {
				prefix = prefix + "/" + info.Name()
			} else if runtime.GOOS == "windows" {
				prefix = prefix + "\\" + info.Name()
			} else {
				prefix = prefix + "/" + info.Name()
			}
		}

		fileInfos, err := file.Readdir(-1)
		if err != nil {
			return fmt.Errorf("os.File.Readdir [%s] failed: %s", file.Name(), err.Error())
		}
		for _, fi := range fileInfos {
			var filename string
			if runtime.GOOS == "linux" {
				filename = file.Name() + "/" + fi.Name()
			} else if runtime.GOOS == "windows" {
				filename = file.Name() + "\\" + fi.Name()
			} else {
				filename = file.Name() + "/" + fi.Name()
			}
			f, err := os.Open(filename)
			if err != nil {
				return fmt.Errorf("os.Open file [%s] failed: %s", filename, err.Error())
			}
			err = compress(f, prefix, zw)
			if err != nil {
				return err
			}
		}
	} else {
		header, err := zip.FileInfoHeader(info)
		if err != nil {
			return err
		}
		if prefix != "" {
			if runtime.GOOS == "linux" {
				header.Name = prefix + "/" + header.Name
			} else if runtime.GOOS == "windows" {
				header.Name = prefix + "\\" + header.Name
			} else {
				header.Name = prefix + "/" + header.Name
			}
		}
		header.Method = zip.Deflate
		fw, err := zw.CreateHeader(header)
		if err != nil {
			return err
		}

		_, err = io.Copy(fw, file)
		if err != nil {
			return err
		}
	}
	return nil
}
