package utility

import (
	"fmt"
	uuid "github.com/satori/go.uuid"
	"keeper/app/pkg/logger"
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

func NewJsonLinesDatabase(filename string) *JsonLinesDatabase {
	return &JsonLinesDatabase{
		Filename: filename,
	}
}

func (j *JsonLinesDatabase) Insert(obj map[string]interface{}) (map[string]interface{}, error) {
	j.ensureLoaded()
	_id, ok := obj["_id"]
	if ok && _id.(string) != "" {
		return nil, fmt.Errorf("Cannot insert duplicate ID %s into %s", _id.(string), j.Filename)
	}

	elem := tools.DeepCopyLooseMap(obj)
	elem["_id"] = uuid.NewV4().String()

	//验证obj的唯一性，除去key字段，所有key对应的值都要一致。
	//unique
	j.Data = append(j.Data, elem)
	if err := j.save(); err != nil {
		logger.Errorf("insert database failed %v", err)
		return nil, err
	}
	return elem, nil
}

func (j *JsonLinesDatabase) Get(id string) map[string]interface{} {
	j.ensureLoaded()
	for _, obj := range j.Data {
		if obj["_id"] != nil && obj["_id"].(string) == id {
			return obj
		}
	}

	return nil
}

//todo
func (j *JsonLinesDatabase) Find() {
	j.ensureLoaded()

}

func (j *JsonLinesDatabase) Update(obj map[string]interface{}) (map[string]interface{}, error) {
	j.ensureLoaded()
	for _, x := range j.Data {
		if obj["_id"] != nil && x["_id"] != nil && x["_id"] == obj["_id"] {
			x = obj
		}
	}

	if err := j.save(); err != nil {
		logger.Errorf("update database failed %v", err)
		return nil, err
	}

	return obj, nil
}

//todo
func (j *JsonLinesDatabase) Path(id string, values ...interface{}) {

}

func (j *JsonLinesDatabase) Remove(id string) (map[string]interface{}, error) {
	j.ensureLoaded()
	var removed map[string]interface{}
	for i, obj := range j.Data {
		if obj["_id"] != nil && obj["_id"].(string) == id {
			removed = obj
			j.Data = append(j.Data[:i], j.Data[i+1:]...) // 删除中间N个元素
		}
	}

	return removed, nil
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

func (j *JsonLinesDatabase) save() error {
	if !j.LoadedOk {
		return fmt.Errorf("not laded")
	}

	if err := tools.WriteFileAllPool(j.Filename, j.Data); err != nil {
		return err
	}

	return nil
}

//
