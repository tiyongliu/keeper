package utility

import (
	"bufio"
	"fmt"
	uuid "github.com/satori/go.uuid"
	"io/ioutil"
	"keeper/app/tools"
	"os"
	"sync"
)

var lock sync.RWMutex

type JsonLinesDatabase struct {
	Filename      string                   `json:"filename"`
	LoadedOk      bool                     `json:"loadedOk"`
	LoadPerformed bool                     `json:"loadPerformed"`
	Data          []map[string]interface{} `json:"data"`
}

/*

 this.data = [];
    this.loadedOk = false;
    this.loadPerformed = false;
*/

func NewJsonLinesDatabase(filename string) *JsonLinesDatabase {
	return &JsonLinesDatabase{
		Filename: filename,
	}
}

func (j *JsonLinesDatabase) Insert(obj map[string]interface{}) {
	/*
	 if (obj._id && (await this.get(obj._id))) {
	      throw new Error(`Cannot insert duplicate ID ${obj._id} into ${this.filename}`);
	    }
	*/
	j.ensureLoaded()
	_id, ok := obj["_id"]
	if ok && _id.(string) != "" {
		fmt.Errorf("Cannot insert duplicate ID %s into %s", _id.(string), j.Filename)
	}

	elem := make(map[string]interface{})
	if err := tools.DeepCopy(elem, obj); err != nil {
		return
	}

	elem["_id"] = uuid.NewV4().String()
	j.Data = append(j.Data, elem)
	j.save()
}

func (j *JsonLinesDatabase) ensureLoaded() {
	if !j.LoadPerformed {
		lock.Lock()
		defer lock.Unlock()

		if !tools.IsExist(j.Filename) {
			j.LoadedOk = true
			j.LoadPerformed = true
			return
		}

		line, err := tools.ReadFileAllPool(j.Filename)
		if err != nil {
			return
		}
		j.Data = line
		j.LoadedOk = true
		j.LoadPerformed = true
	}
}

func (j *JsonLinesDatabase) save() {
	if !j.LoadedOk {
		return
	}
	//写文件

}

func WriteToFile(fileName string, content string) error {
	f, err := os.OpenFile(fileName, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0644)
	if err != nil {
		fmt.Println("file create failed. err: " + err.Error())
	} else {
		// offset
		//os.Truncate(filename, 0) //clear
		n, _ := f.Seek(0, os.SEEK_END)
		_, err = f.WriteAt([]byte(content), n)
		fmt.Println("write succeed!")
		defer f.Close()
	}
	return err
}

//将map[string]string写入文件
func WriteMaptoFile(m map[string]string, filePath string) error {
	f, err := os.Create(filePath)
	if err != nil {
		fmt.Printf("create map file error: %v\n", err)
		return err
	}
	defer f.Close()

	w := bufio.NewWriter(f)
	for k, v := range m {
		lineStr := fmt.Sprintf("%s^%s", k, v)
		fmt.Fprintln(w, lineStr)
	}
	return w.Flush()
}

//这种方式每次都会覆盖 test.txt内容，如果test.txt文件不存在会创建。
func writeMap() {
	content := []byte("测试1\n测试2\n")
	err := ioutil.WriteFile("test.txt", content, 0644)
	if err != nil {
		panic(err)
	}
}
