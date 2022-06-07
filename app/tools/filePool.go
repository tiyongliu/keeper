package tools

import (
	"bufio"
	"io"
	"os"
	"strings"
)

//读取所有文件读连接池
func ReadFileAllPool(name string) ([]map[string]interface{}, error) {
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
