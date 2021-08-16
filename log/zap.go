package logs

import (
	"io"

	"github.com/labstack/gommon/log"
	"go.uber.org/zap"
)

//go:generate moq -out zap_mock_test.go . ZapInterface

type Zap struct {
	*zap.SugaredLogger
	*zap.AtomicLevel
}

func (z *Zap) Printf(format string, args ...interface{}) {
	z.Infof(format, args...)
}

func (z *Zap) Print(args ...interface{}) {
	z.Info(args...)
}

func (z *Zap) Output() io.Writer {
	return z.Output() // Info? or Open
}

func (z *Zap) SetOutput(w io.Writer) {
	z.SetOutput(w) //ErrorOutput? or NewStdLog
}

func (z *Zap) Prefix() string {
	return z.Prefix()
}

func (z *Zap) SetPrefix(p string) {
	z.SetPrefix(p)
}

func (z *Zap) Level() log.Lvl {
	return log.Lvl(z.AtomicLevel.Level()) // AtomicLevel - Level()
}

func (z *Zap) SetLevel(v log.Lvl) {
	z.SetLevel(v) // AtomicLevel - SetLevel()
}

func (z *Zap) SetHeader(h string) {
	z.SetHeader(h) // Not Sure what header refers to
}

func (z *Zap) Printj(j log.JSON) {
	z.Info(j)
}

func (z *Zap) Debugj(j log.JSON) {
	z.Debug(j)
}

func (z *Zap) Infoj(j log.JSON) {
	z.Info(j)
}

func (z *Zap) Warnj(j log.JSON) {
	z.Warn(j)
}

func (z *Zap) Errorj(j log.JSON) {
	z.Error(j)
}

func (z *Zap) Fatalj(j log.JSON) {
	z.Fatal(j)
}

func (z *Zap) Panicj(j log.JSON) {
	z.Panic(j)
}
