package tools

import (
	"bufio"
	"io"
	"io/ioutil"
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
	//list, err := os.OpenFile(name, os.O_RDWR|os.O_TRUNC|os.O_CREATE, 0644)
	//if err != nil {
	//	return
	//}
	//
	//for _, item := range dataSource {
	//
	//}

	content, err := JsonMarshal(dataSource)
	if err != nil {
		return err
	}

	return ioutil.WriteFile("test.txt", content, 0644)
}
