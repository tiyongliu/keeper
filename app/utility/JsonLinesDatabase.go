package utility

import (
	"fmt"
	uuid "github.com/satori/go.uuid"
	"keeper/app/tools"
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
