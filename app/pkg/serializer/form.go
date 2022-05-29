package serializer

import (
	"errors"
	"github.com/astaxie/beego/validation"
	"github.com/gin-gonic/gin"
	"keeper/app/pkg/logger"
	"strings"
)

var (
	InvalidParams = errors.New("invalid params")
)

func Verify(c *gin.Context, form interface{}) error {
	if err := c.Bind(form); err != nil {
		logger.Errorf("request failed: url %s, err %v", c.Request.URL, err)
		return InvalidParams
	}
	valid := validation.Validation{}
	check, err := valid.Valid(form)
	if err != nil {
		return err
	}
	if !check {
		markErrors(valid.Errors)
		return buildFormErr(valid.Errors)
	}
	return nil
}

func markErrors(errors []*validation.Error) {
	for _, err := range errors {
		logger.Errorf("validate rules: label %s, message %s", err.Key, err.Message)
	}
	return
}

func buildFormErr(errs []*validation.Error) error {
	var msg strings.Builder
	for _, v := range errs {
		msg.WriteString(v.Error())
	}
	return errors.New(msg.String())
}
