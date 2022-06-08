package serializer

import (
	"context"
	"time"
)

// Response 基础序列化器
type Response struct {
	Code    int         `json:"code"`
	Result  interface{} `json:"result"`
	Message string      `json:"message"`
	Type    string      `json:"type"`
	Time    int64       `json:"time"`
}

const (
	Code_SUCCESS = iota
	Code_ERROR
)

func SuccessData(ctx context.Context, message string, data interface{}) *Response {
	return &Response{
		Code:    Code_SUCCESS,
		Result:  data,
		Message: message,
		Type:    "success",
		Time:    time.Now().Unix(),
	}
}

func Fail(ctx context.Context, message string) *Response {
	return &Response{
		Code:    Code_ERROR,
		Result:  nil,
		Message: message,
		Time:    time.Now().Unix(),
	}
}
