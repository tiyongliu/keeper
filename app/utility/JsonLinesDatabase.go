package utility

import (
	"fmt"
	uuid "github.com/satori/go.uuid"
	"io/ioutil"
	"keeper/app/pkg/logger"
	"keeper/app/tools"
	"sync"
)

var lock sync.RWMutex

const database_key = "_id"

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
	j.EnsureLoaded()
	dynamicId, ok := obj[database_key]
	if ok && dynamicId.(string) != "" {
		return nil, fmt.Errorf("Cannot insert duplicate ID %s into %s", dynamicId.(string), j.Filename)
	}

	elem := tools.DeepCopyUnknownMap(obj)
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
	j.EnsureLoaded()
	for _, obj := range j.Data {
		if obj[database_key] != nil && obj[database_key].(string) == id {
			return obj
		}
	}

	return nil
}

func (j *JsonLinesDatabase) Find() []map[string]interface{} {
	j.EnsureLoaded()
	return j.Data
}

func (j *JsonLinesDatabase) Update(obj map[string]interface{}) (map[string]interface{}, error) {
	j.EnsureLoaded()
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

//todo
func (j *JsonLinesDatabase) Path(id string, values ...interface{}) {

}

func (j *JsonLinesDatabase) Remove(id string) (map[string]interface{}, error) {
	j.EnsureLoaded()
	var removed map[string]interface{}
	for i, obj := range j.Data {
		if obj[database_key] != nil && obj[database_key].(string) == id {
			removed = obj
			j.Data = append(j.Data[:i], j.Data[i+1:]...) // 删除中间N个元素
		}
	}

	return removed, nil
}

func (j *JsonLinesDatabase) EnsureLoaded() {
	fmt.Printf("LoadPerformed: ----------: %v", j.LoadPerformed)
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
			fmt.Printf("err110: ----------: %s\n", err)

			return
		}
		j.Data = line
		j.LoadedOk = true
		j.LoadPerformed = true

		f, err := ioutil.ReadFile(j.Filename)
		if err != nil {
			fmt.Println("read fail", err)
		}

		fmt.Printf("string(f)string(f): ----------: %s", string(f))
	}
}

func (j *JsonLinesDatabase) save() error {
	if !j.LoadedOk {
		return fmt.Errorf("not laded")
	}

	lock.Lock()
	defer lock.Unlock()
	if err := tools.WriteFileAllPool(j.Filename, j.Data); err != nil {
		return err
	}

	return nil
}

//
