package utility

import (
	"github.com/samber/lo"
	"os"
)

func IsExist(path string) bool {
	_, err := os.Stat(path)
	if err != nil {
		if os.IsExist(err) {
			return true
		}
		if os.IsNotExist(err) {
			return false
		}
		return false
	}
	return true
}

func MapOmit(object map[string]string, paths []string) map[string]string {
	for key := range object {
		if lo.Contains[string](paths, key) {
			delete(object, key)
		}
	}

	return object
}

func MapValues(object map[string]string, iteratee func(k, v interface{}) interface{}) map[string]string {
	result := map[string]string{}
	for key, value := range object {
		v := iteratee(key, value)
		if v != nil {
			result[key] = v.(string)
		}
	}

	return result
}

func FilterUnknownMap(object map[string]interface{}, keys ...string) map[string]interface{} {
	result := map[string]interface{}{}
	for key, value := range object {
		if !lo.Contains[string](keys, key) {
			result[key] = value
		}
	}

	return result
}
