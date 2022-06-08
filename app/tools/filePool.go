package tools

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
)

//读取所有文件读连接池
func ReadFileAllPool1(name string) ([]map[string]interface{}, error) {
	fd, err := os.Open(name)
	defer fd.Close()
	if err != nil {
		return nil, err
	}
	buff := bufio.NewReader(fd)

	var list []map[string]interface{}
	for {
		data, _, eof := buff.ReadLine()
		if eof == io.EOF {
			break
		}
		text := strings.TrimSpace(string(data))
		if unmarshal, err := JsonUnmarshal([]byte(text)); err != nil {
			fmt.Printf("err: ----------: %v\n", err)
			break
		} else {
			list = append(list, unmarshal)
		}
	}

	return list, nil
}

func ReadFileAllPool(name string) ([]map[string]interface{}, error) {
	file, err := os.Open(name)
	if err != nil {
		fmt.Println("open file err=", err)
		return nil, nil
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
		fmt.Printf("bufio.NewReader(file): ----------: %s\n", str)
		if err == io.EOF { // io.EOF表示文件的末尾
			break
		}
		text := strings.TrimSpace(str)
		if unmarshal, err := JsonUnmarshal([]byte(text)); err != nil {
			fmt.Printf("err: ----------: %v\n", err)
			break
		} else {
			list = append(list, unmarshal)
		}
	}
	return list, nil
}

func WriteFileAllPool(name string, dataSource []map[string]interface{}) error {
	f, err := os.OpenFile(name, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	n, _ := f.Seek(0, os.SEEK_END)
	if content, err := JsonMarshal(dataSource); err == nil {
		_, err = f.WriteAt(content, n)
	}
	defer f.Close()
	return err
}
