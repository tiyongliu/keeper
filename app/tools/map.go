package tools

import (
	"encoding/json"
	"fmt"
)

func MapBoundString(m map[string]interface{}) map[string]string {
	ret := make(map[string]string, len(m))
	for k, v := range m {
		ret[k] = fmt.Sprint(v)
	}
	return ret
}

func LooseMapValue(m map[string]string) map[string]interface{} {
	nm := make(map[string]interface{})
	for k, v := range m {
		nm[k] = v
	}

	return nm
}

func DeepCopyLooseMap(valueMap map[string]interface{}) map[string]interface{} {
	newMap := make(map[string]interface{})
	for k, v := range valueMap {
		newMap[k] = v
	}
	return newMap
}

func UniqueMap(list []map[string]interface{}, valueMap map[string]interface{}) bool {
	for _, item := range list {
		CompareTwoMapInterface(item, valueMap)
	}

	return true
}

func CompareTwoMapInterface(data1 map[string]interface{}, data2 map[string]interface{}) bool {
	keySlice := make([]string, 0)
	dataSlice1 := make([]interface{}, 0)
	dataSlice2 := make([]interface{}, 0)
	for key, value := range data1 {
		keySlice = append(keySlice, key)
		dataSlice1 = append(dataSlice1, value)
	}
	for _, key := range keySlice {
		if data, ok := data2[key]; ok {
			dataSlice2 = append(dataSlice2, data)
		} else {
			return false
		}
	}
	dataStr1, _ := json.Marshal(dataSlice1)
	dataStr2, _ := json.Marshal(dataSlice2)

	return string(dataStr1) == string(dataStr2)
}
