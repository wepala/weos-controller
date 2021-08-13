package logs

import (
	"testing"

	"go.uber.org/zap"
)

func TestZap_Printf(t *testing.T) {
	zlogger, _ := zap.NewProduction()
	defer zlogger.Sync() // flushes buffer, if any
	sugar := zlogger.Sugar()
	logger := Zap{sugar}
	logger.Printf("test")
}

func TestZap_Prefix(t *testing.T) {
	zlogger, _ := zap.NewProduction()
	defer zlogger.Sync() // flushes buffer, if any
	sugar := zlogger.Sugar()
	logger := Zap{sugar}
	logger.Prefix()
}

func TestZap_SetPrefix(t *testing.T) {
	zlogger, _ := zap.NewProduction()
	defer zlogger.Sync() // flushes buffer, if any
	sugar := zlogger.Sugar()
	logger := Zap{sugar}
	prefix := "Urgent"
	logger.SetPrefix(prefix)
}

func TestZap_Level(t *testing.T) {
	zlogger, _ := zap.NewProduction()
	defer zlogger.Sync() // flushes buffer, if any
	sugar := zlogger.Sugar()
	logger := Zap{sugar}
	logger.Level()
}

func TestZap_SetLevel(t *testing.T) {
	zlogger, _ := zap.NewProduction()
	defer zlogger.Sync() // flushes buffer, if any
	sugar := zlogger.Sugar()
	logger := Zap{sugar}
	logger.SetLevel(logger.Level())
}

func TestZap_Output(t *testing.T) {
	zlogger, _ := zap.NewProduction()
	defer zlogger.Sync() // flushes buffer, if any
	sugar := zlogger.Sugar()
	logger := Zap{sugar}
	logger.Output()
}

func TestZap_SetOutput(t *testing.T) {
	zlogger, _ := zap.NewProduction()
	defer zlogger.Sync() // flushes buffer, if any
	sugar := zlogger.Sugar()
	logger := Zap{sugar}
	logger.SetOutput(logger.Output())
}

func TestZap_SetHeader(t *testing.T) {
	zlogger, _ := zap.NewProduction()
	defer zlogger.Sync() // flushes buffer, if any
	sugar := zlogger.Sugar()
	logger := Zap{sugar}
	header := "Test Header"
	logger.SetHeader(header)
}

func TestZap_Debug(t *testing.T) {
	zlogger, _ := zap.NewProduction()
	defer zlogger.Sync() // flushes buffer, if any
	sugar := zlogger.Sugar()
	logger := Zap{sugar}

	Debugged := "The issue is xyz"

	logger.Debug(Debugged)
}
