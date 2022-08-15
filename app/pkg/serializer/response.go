package serializer

import (
	"time"
)

// Response 基础序列化器
type Response struct {
	Status  int         `json:"status"`
	Result  interface{} `json:"result"`
	Message string      `json:"message"`
	Type    string      `json:"type"`
	Time    int64       `json:"time"`
}

const (
	StatusCodeSuccess = iota
	StatusCodeFailed
)

func SuccessData(message string, data interface{}) *Response {
	return &Response{
		Status:  StatusCodeSuccess,
		Result:  data,
		Message: message,
		Time:    time.Now().Unix(),
	}
}

func Fail(message string) *Response {
	return &Response{
		Status:  StatusCodeFailed,
		Result:  nil,
		Message: message,
		Time:    time.Now().Unix(),
	}
}
