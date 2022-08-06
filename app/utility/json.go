package utility

import "encoding/json"

func ToJsonStr(o interface{}) string {
	jsonBytes, err := json.Marshal(o)
	if err != nil {
		return err.Error()
	}
	return string(jsonBytes)
}

//解析json
func JsonUnmarshal(data []byte) (map[string]interface{}, error) {
	msg := make(map[string]interface{})
	err := json.Unmarshal(data, &msg)
	if err != nil {
		return nil, err
	}
	return msg, nil
}

//编码json
func JsonMarshal(data interface{}) ([]byte, error) {
	msg, err := json.Marshal(&data)
	if err != nil {
		return nil, err
	}
	return msg, nil
}
