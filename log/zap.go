package logs

import (
	"io"

	"github.com/labstack/gommon/log"

	"go.uber.org/zap"
)

type Zap struct {
	*zap.SugaredLogger
}

func (z *Zap) Printf(format string, args ...interface{}) {
	log.Infof(format, args...)
}

func (z *Zap) Print(args ...interface{}) {
	log.Info(args...)
}

func (z *Zap) Output() io.Writer {
	return log.Output()
}

func (z *Zap) SetOutput(w io.Writer) {
	log.SetOutput(w)
}

func (z *Zap) Prefix() string {
	return log.Prefix()
}

func (z *Zap) SetPrefix(p string) {
	log.SetPrefix(p)
}

func (z *Zap) Level() log.Lvl {
	return log.Level()
}

func (z *Zap) SetLevel(v log.Lvl) {
	log.SetLevel(v)
}

func (z *Zap) SetHeader(h string) {
	log.SetHeader(h)
}

func (z *Zap) Printj(j log.JSON) {
	log.Infoj(j)
}

func (z *Zap) Debugj(j log.JSON) {
	log.Debugj(j)
}

func (z *Zap) Infoj(j log.JSON) {
	log.Infoj(j)
}

func (z *Zap) Warnj(j log.JSON) {
	log.Warnj(j)
}

func (z *Zap) Errorj(j log.JSON) {
	log.Errorj(j)
}

func (z *Zap) Fatalj(j log.JSON) {
	log.Fatalj(j)
}

func (z *Zap) Panicj(j log.JSON) {
	log.Panicj(j)
}
