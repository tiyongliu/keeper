package utility

import (
	"bufio"
	"io"
	"os"
	"strings"
)

//读取所有文件读连接池
func ReadFileAllPool(name string) ([]map[string]interface{}, error) {
	file, err := os.Open(name)
	if err != nil {
		return nil, err
	}

	//当函数退出时，要及时的关闭file
	defer file.Close() //要及时关闭file句柄，否则会有内存泄漏.

	// 创建一个 *Reader  ，是带缓冲的
	/*
		const (
		defaultBufSize = 4096 //默认的缓冲区为4096
		)
	*/
	var list []map[string]interface{}
	reader := bufio.NewReader(file)

	//循环的读取文件的内容
	for {
		str, err := reader.ReadString('\n') // 读到一个换行就结束
		if err == io.EOF {                  // io.EOF表示文件的末尾
			break
		}
		text := strings.TrimSpace(str)
		if unmarshal, err := JsonUnmarshal([]byte(text)); err != nil {
			break
		} else {
			list = append(list, unmarshal)
		}
	}
	return list, nil
}

func WriteFileAllPool(name string, dataSource []map[string]interface{}) error {
	file, err := os.OpenFile(name, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, os.ModePerm)
	if err != nil {
		return err
	}

	defer file.Close()
	//写入文件时，使用带缓存的 *Writer
	write := bufio.NewWriter(file)
	for _, x := range dataSource {
		if marshal, err := JsonMarshal(x); err == nil {
			write.WriteString(string(marshal) + "\n")
		}
	}

	//Flush将缓存的文件真正写入到文件中
	return write.Flush()
}
