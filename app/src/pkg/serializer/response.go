package serializer

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

const (
	Code_SUCCESS = iota
	Code_ERROR
)

// Response 基础序列化器
type Response struct {
	Code    int         `json:"code"`
	Result  interface{} `json:"result,omitempty"`
	Message string      `json:"message"`
	Type    string      `json:"type"`
	Time    int64       `json:"time"`
}

//失败
func Fail(ctx *gin.Context, message string) {
	ctx.JSON(http.StatusOK, &Response{
		Code:    Code_ERROR,
		Result:  nil,
		Message: message,
		Time:    time.Now().Unix(),
	})
}

//成功
func Success(ctx *gin.Context, message string) {
	ctx.JSON(http.StatusOK, &Response{
		Code:    Code_SUCCESS,
		Result:  nil,
		Message: message,
		Type:    "success",
		Time:    time.Now().Unix(),
	})
}

func SuccessData(ctx *gin.Context, message string, data interface{}) {
	ctx.JSON(http.StatusOK, &Response{
		Code:    Code_SUCCESS,
		Result:  data,
		Message: message,
		Type:    "success",
		Time:    time.Now().Unix(),
	})
}

func Reply(ctx *gin.Context, res *Response) {
	res.Type = "success"
	res.Time = time.Now().Unix()
	ctx.JSON(http.StatusOK, res)
}

func Echo(ctx *gin.Context, code int, message string, data interface{}) {
	ctx.JSON(http.StatusOK, &Response{
		Code:    code,
		Result:  data,
		Message: message,
		Type:    "success",
		Time:    time.Now().Unix(),
	})
}
