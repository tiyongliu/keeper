package tools

import "encoding/json"

func ToJsonStr(o interface{}) string {
	jsonBytes, err := json.Marshal(o)
	if err != nil {
		return err.Error()
	}
	return string(jsonBytes)
}
