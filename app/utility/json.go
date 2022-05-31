package utility

import "encoding/json"

//解析json
func JSONUnmarshal(data []byte) (map[string]interface{}, error) {
	msg := make(map[string]interface{})
	err := json.Unmarshal(data, &msg)
	if err != nil {
		return nil, err
	}
	return msg, nil
}

//编码json
func JSONMarshal(data interface{}) ([]byte, error) {
	msg, err := json.Marshal(&data)
	if err != nil {
		return nil, err
	}
	return msg, nil
}
