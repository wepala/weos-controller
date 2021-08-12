package logs

import (
	"github.com/labstack/gommon/log"
	"go.uber.org/zap"
	"io"
)

type Zap struct {
	*zap.SugaredLogger
}

func (z *Zap) Printf(format string, args ...interface{}) {
	panic("implement me")
}

func (z *Zap) Print(args ...interface{}) {
	panic("implement me")
}

func (z *Zap) Output() io.Writer {
	panic("implement me")
}

func (z *Zap) SetOutput(w io.Writer) {
	panic("implement me")
}

func (z *Zap) Prefix() string {
	panic("implement me")
}

func (z *Zap) SetPrefix(p string) {
	panic("implement me")
}

func (z *Zap) Level() log.Lvl {
	panic("implement me")
}

func (z *Zap) SetLevel(v log.Lvl) {

}

func (z Zap) SetHeader(h string) {
	panic("implement me")
}

func (z Zap) Printj(j log.JSON) {
	panic("implement me")
}

func (z Zap) Debugj(j log.JSON) {
	panic("implement me")
}

func (z Zap) Infoj(j log.JSON) {
	panic("implement me")
}

func (z Zap) Warnj(j log.JSON) {
	panic("implement me")
}

func (z Zap) Errorj(j log.JSON) {
	panic("implement me")
}

func (z Zap) Fatalj(j log.JSON) {
	panic("implement me")
}

func (z Zap) Panicj(j log.JSON) {
	panic("implement me")
}
