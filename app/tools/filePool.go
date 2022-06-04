package tools

import (
	"bufio"
	"io"
	"os"
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
		if unmarshal, err := JsonUnmarshal(data); err != nil {
			break
		} else {
			list = append(list, unmarshal)
		}
	}

	return list, nil
}

func WriteFileAllPool() {

}
