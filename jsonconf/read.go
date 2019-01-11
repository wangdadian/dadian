package jsonconf

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
)

func init() {

}

// 读取szFile文件的json配置文件
// 去除以“#”为首的行（自定义注释行，实际是不允许的）
// 清空空格以及换行符
// 返回：正常情况返回字节流数据以及nil，失败返回空以及相应错误信息
func ReadConf(szFile string) ([]byte, error) {
	file, err := os.Open(szFile)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	bufro := bufio.NewReader(file)
	var szLine string
	var byRetBuff []byte = nil
	bEOF := false
	for {
		if bEOF {
			break
		}
		szLine, err = bufro.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				// 文件读完毕，继续处理已读取的内容
				bEOF = true
			} else {
				// 其他错误
				return nil, err
			}
		}
		// 删除行首尾空格/TAB
		for {
			if strings.HasPrefix(szLine, " ") || strings.HasPrefix(szLine, "\t") ||
				strings.HasSuffix(szLine, " ") || strings.HasSuffix(szLine, "\t") {
				szLine = strings.Trim(szLine, " ")
				szLine = strings.Trim(szLine, "\t")
			} else {
				break
			}
		}
		// 删除行末的换行符
		szLine = strings.TrimRight(szLine, "\n")
		szLine = strings.TrimRight(szLine, "\r")
		// 注释行，不处理
		if strings.HasPrefix(szLine, "#") || len(szLine) == 0 {
			continue
		}
		byRetBuff = append(byRetBuff, []byte(szLine)...)
	}
	// 验证json是否合法
	var r interface{}
	err = json.Unmarshal(byRetBuff, &r)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("invalid json data from %s: %s", szFile, err.Error()))
	}

	return byRetBuff, nil
}
