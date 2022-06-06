package tools

import "fmt"

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
