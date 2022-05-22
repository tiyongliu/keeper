package main

import (
	"fmt"
	"reflect"
)

//reflect.TypeOf(connection).Kind().String()
func main() {
	r := map[string]string{"v": "1"}
	d := reflect.TypeOf(r).Kind() == reflect.Map

	fmt.Println(d)
}
