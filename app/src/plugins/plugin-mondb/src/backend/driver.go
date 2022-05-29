package backend

import (
	"context"
	"encoding/json"
	"fmt"
	"keeper/app/src/modules"
	"keeper/app/src/pkg/serializer"
	plugin_mondb "keeper/app/src/plugins/plugin-mondb"
	"time"
)

type MMMM struct {
	ctx context.Context
}

func NewMMMM() *MMMM {
	return &MMMM{}
}

func (m *MMMM) GetVersion(params map[string]interface{}) (*serializer.Response, error) {
	if params != nil {
		form := modules.SimpleSettingMongoDB{}
		arr, err := json.Marshal(form)
		if err != nil {
			return nil, err
		}

		err = json.Unmarshal(arr, &form)
		if err != nil {
			fmt.Println(err)
			return nil, err
		}

		fmt.Println(form)

		pool, err := plugin_mondb.NewSimpleMongoDBPool(&modules.SimpleSettingMongoDB{
			Host: "localhost",
			Port: "27017",
		})

		defer func() {
			if err == nil {
				pool.Close()
			}
		}()

		if err != nil {
			return &serializer.Response{
				Code:    serializer.Code_ERROR,
				Result:  nil,
				Message: err.Error(),
				Type:    "error",
				Time:    time.Now().Unix(),
			}, err
		}

		getVersion, err := pool.GetVersion()
		if err != nil {
			return &serializer.Response{
				Code:    serializer.Code_ERROR,
				Result:  nil,
				Message: err.Error(),
				Type:    "error",
				Time:    time.Now().Unix(),
			}, err
		}

		return &serializer.Response{
			Code:    serializer.Code_SUCCESS,
			Result:  getVersion,
			Message: "",
			Type:    "success",
			Time:    time.Now().Unix(),
		}, err

		return nil, nil
	}
	//form, ok := ctx.Value(q).(modules.SimpleSettingMongoDB)
	//if !ok {
	//	return &serializer.Response{
	//		Code:    serializer.Code_ERROR,
	//		Result:  nil,
	//		Message: "",
	//		Type:    "error",
	//		Time:    time.Now().Unix(),
	//	}, nil
	//}
	//
	//pool, err := plugin_mondb.NewSimpleMongoDBPool(&form)
	//
	//defer func() {
	//	if err == nil {
	//		pool.Close()
	//	}
	//}()
	//
	//if err != nil {
	//	return &serializer.Response{
	//		Code:    serializer.Code_ERROR,
	//		Result:  nil,
	//		Message: err.Error(),
	//		Type:    "error",
	//		Time:    time.Now().Unix(),
	//	}, err
	//}
	//
	//getVersion, err := pool.GetVersion()
	//if err != nil {
	//	return &serializer.Response{
	//		Code:    serializer.Code_ERROR,
	//		Result:  nil,
	//		Message: err.Error(),
	//		Type:    "error",
	//		Time:    time.Now().Unix(),
	//	}, err
	//}
	//
	//return &serializer.Response{
	//	Code:    serializer.Code_SUCCESS,
	//	Result:  getVersion,
	//	Message: "",
	//	Type:    "success",
	//	Time:    time.Now().Unix(),
	//}, err

	return nil, nil
}
