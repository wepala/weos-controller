package logs

import (
	"go.uber.org/zap"
	"testing"
)

func TestZap_Print(t *testing.T) {
	zlogger, _ := zap.NewProduction()
	defer zlogger.Sync() // flushes buffer, if any
	sugar := zlogger.Sugar()
	logger := Zap{sugar}
	logger.Print("test")
}
