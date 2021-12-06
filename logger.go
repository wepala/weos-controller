package weoscontroller

import (
	log "github.com/sirupsen/logrus"
	weosLogs "github.com/wepala/weos-controller/log"
)

//LoggerWrapper makes a WeOS compatible logger
type LoggerWrapper struct {
	logger weosLogs.Zap
}

func (l LoggerWrapper) WithField(key string, value interface{}) *log.Entry {
	panic("implement me")
}

func (l LoggerWrapper) WithFields(fields log.Fields) *log.Entry {
	panic("implement me")
}

func (l LoggerWrapper) WithError(err error) *log.Entry {
	panic("implement me")
}

func (l LoggerWrapper) Debugf(format string, args ...interface{}) {
	l.logger.Debugf(format, args...)
}

func (l LoggerWrapper) Infof(format string, args ...interface{}) {
	l.logger.Infof(format, args...)
}

func (l LoggerWrapper) Printf(format string, args ...interface{}) {
	l.logger.Printf(format, args...)
}

func (l LoggerWrapper) Warnf(format string, args ...interface{}) {
	l.logger.Warnf(format, args...)
}

func (l LoggerWrapper) Warningf(format string, args ...interface{}) {
	l.logger.Warnf(format, args...)
}

func (l LoggerWrapper) Errorf(format string, args ...interface{}) {
	l.logger.Errorf(format, args...)
}

func (l LoggerWrapper) Fatalf(format string, args ...interface{}) {
	l.logger.Fatalf(format, args...)
}

func (l LoggerWrapper) Panicf(format string, args ...interface{}) {
	l.logger.Panicf(format, args...)
}

func (l LoggerWrapper) Debug(args ...interface{}) {
	l.logger.Debug(args...)
}

func (l LoggerWrapper) Info(args ...interface{}) {
	l.logger.Info(args...)
}

func (l LoggerWrapper) Print(args ...interface{}) {
	l.logger.Print(args...)
}

func (l LoggerWrapper) Warn(args ...interface{}) {
	l.logger.Warn(args...)
}

func (l LoggerWrapper) Warning(args ...interface{}) {
	l.logger.Warn(args...)
}

func (l LoggerWrapper) Error(args ...interface{}) {
	l.logger.Error(args...)
}

func (l LoggerWrapper) Fatal(args ...interface{}) {
	l.logger.Fatal(args...)
}

func (l LoggerWrapper) Panic(args ...interface{}) {
	l.logger.Panic(args...)
}

func (l LoggerWrapper) Debugln(args ...interface{}) {
	l.logger.Debug(args...)
}

func (l LoggerWrapper) Infoln(args ...interface{}) {
	l.logger.Info(args...)
}

func (l LoggerWrapper) Println(args ...interface{}) {
	l.logger.Print(args...)
}

func (l LoggerWrapper) Warnln(args ...interface{}) {
	l.logger.Warn(args...)
}

func (l LoggerWrapper) Warningln(args ...interface{}) {
	l.logger.Warn(args...)
}

func (l LoggerWrapper) Errorln(args ...interface{}) {
	l.logger.Error(args...)
}

func (l LoggerWrapper) Fatalln(args ...interface{}) {
	panic("implement me")
}

func (l LoggerWrapper) Panicln(args ...interface{}) {
	panic("implement me")
}

func (l LoggerWrapper) Tracef(format string, args ...interface{}) {
	l.logger.Debugf(format, args...)
}

func (l LoggerWrapper) Trace(args ...interface{}) {
	l.logger.Debug(args...)
}

func (l LoggerWrapper) Traceln(args ...interface{}) {
	panic("implement me")
}

func NewLogger(l weosLogs.Zap) log.Ext1FieldLogger {
	return &LoggerWrapper{
		logger: l,
	}
}
