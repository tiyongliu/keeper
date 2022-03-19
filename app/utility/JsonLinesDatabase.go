package utility

import (
	"errors"
	"fmt"
	uuid "github.com/satori/go.uuid"
	"keeper/app/pkg/logger"
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
	dynamicId, ok := obj[database_key]
	if ok && dynamicId.(string) != "" {
		return nil, fmt.Errorf("cannot insert duplicate ID %s into %s", dynamicId.(string), j.Filename)
	}

	elem := DeepCopyUnknownMap(obj)
	elem[database_key] = uuid.NewV4().String()
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
		if obj[database_key] != nil && obj[database_key].(string) == id {
			return obj
		}
	}

	return nil
}

func (j *JsonLinesDatabase) Find() []map[string]interface{} {
	j.ensureLoaded()
	return j.Data
}

func (j *JsonLinesDatabase) Update(obj map[string]interface{}) (map[string]interface{}, error) {
	j.ensureLoaded()
	for _, x := range j.Data {
		if obj[database_key] != nil && x[database_key] != nil && x[database_key] == obj[database_key] {
			x = obj
		}
	}

	if err := j.save(); err != nil {
		logger.Errorf("update database failed %v", err)
		return nil, err
	}

	return obj, nil
}

func (j *JsonLinesDatabase) Path(id string, values ...interface{}) {

}

func (j *JsonLinesDatabase) Remove(id string) (map[string]interface{}, error) {
	j.ensureLoaded()
	var removed map[string]interface{}
	var match bool
	var err error
	for i, obj := range j.Data {
		if obj[database_key] != nil && obj[database_key].(string) == id {
			match = true
			removed = obj
			j.Data = append(j.Data[:i], j.Data[i+1:]...) // 删除中间N个元素
		}
	}

	if !match {
		err = errors.New("id in not a valid")
		return nil, err
	}
	return removed, j.save()
}

func (j *JsonLinesDatabase) ensureLoaded() {
	if !j.LoadPerformed {
		lock.Lock()
		defer lock.Unlock()

		if !IsExist(j.Filename) {
			j.LoadedOk = true
			j.LoadPerformed = true
			return
		}

		line, err := ReadFileAllPool(j.Filename)
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
	lock.Lock()
	defer lock.Unlock()
	return WriteFileAllPool(j.Filename, j.Data)
}

func (j *JsonLinesDatabase) EnsureOpened(conid string) {
	if conid == "" {
		return
	}
}

func (j *JsonLinesDatabase) EnsureOpened(conid string) {
	if conid == "" {
		return
	}
}
