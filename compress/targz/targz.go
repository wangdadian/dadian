package targz

import (
	"archive/tar"
	"compress/gzip"
	"errors"
	"fmt"
	"io"
	"os"
	"runtime"
	"strings"
)

type TargzCompress struct {
}

func (self *TargzCompress) Compress(szFiles []string, szDestFile string) (string, error) {
	// 目标文件非目录
	if strings.HasSuffix(szDestFile, "/") || strings.HasSuffix(szDestFile, "\\") {
		return "", errors.New("dest file name is a dir")
	}
	if yes := strings.HasSuffix(szDestFile, ".tar.gz"); !yes {
		szDestFile += ".tar.gz"
	}
	// 创建目标文件
	df, err := os.Create(szDestFile)
	if err != nil {
		return "", errors.New(fmt.Sprintf("create dest file[%s] failed: %s", szDestFile, err.Error()))
	}
	defer df.Close()
	gw := gzip.NewWriter(df)
	defer gw.Close()
	tw := tar.NewWriter(gw)
	defer tw.Close()
	// 打开文件，并压缩
	for _, szFile := range szFiles {
		f, err := os.Open(szFile)
		if err != nil {
			err = errors.New(fmt.Sprintf("open file[%s] failed: %s", szFile, err.Error()))
			goto goto_ret
		}
		if err := compress(f, "", tw); err != nil {
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
func (self *TargzCompress) DeCompress(tarFile, destDir string) error {
	srcFile, err := os.Open(tarFile)
	if err != nil {
		return fmt.Errorf("os.Open file[%s] failed: %s", tarFile, err.Error())
	}
	defer srcFile.Close()
	gr, err := gzip.NewReader(srcFile)
	if err != nil {
		return fmt.Errorf("gzip.NewReader failed: %s", err.Error())
	}
	defer gr.Close()
	tr := tar.NewReader(gr)
	for {
		hdr, err := tr.Next()
		if err != nil {
			if err == io.EOF {
				break
			} else {
				return fmt.Errorf("tar.Next error: %s", err.Error())
			}
		}
		filename := destDir
		if runtime.GOOS == "linux" {
			if yes := strings.HasSuffix(destDir, "/"); !yes {
				filename += "/"
			}
		} else if runtime.GOOS == "windows" {
			if yes := strings.HasSuffix(destDir, "\\"); !yes {
				filename += "\\"
			}
		} else {
			if yes := strings.HasSuffix(destDir, "/"); !yes {
				filename += "/"
			}
		}
		filename += hdr.Name
		file, err := createFile(filename)
		if err != nil {
			return err
		}
		// fmt.Printf("create file: %s\n", filename)
		io.Copy(file, tr)
		if err := file.Close(); err != nil {
			// fmt.Printf("close file[%s] failed: %s", file.Name(), err.Error())
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
func compress(file *os.File, prefix string, tw *tar.Writer) error {
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
			err = compress(f, prefix, tw)
			if err != nil {
				return err
			}
		}
	} else {
		header, err := tar.FileInfoHeader(info, "")
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

		err = tw.WriteHeader(header)
		if err != nil {
			return fmt.Errorf("tar.WriteHeader failed: %s", err.Error())
		}
		_, err = io.Copy(tw, file)
		if err != nil {
			return fmt.Errorf("io.Copy failed: %s", err.Error())
		}
	}
	return nil
}
