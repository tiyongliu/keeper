package logger_test

import (
	"keeper/app/pkg/logger"
	"testing"
)

func TestInitLogger(t *testing.T) {
	logger.Infof("test")
}
