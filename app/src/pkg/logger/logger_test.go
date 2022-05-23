package logger_test

import (
	"system/pkg/logger"
	"testing"
)

func TestInitLogger(t *testing.T) {
	logger.Info("test")
}
