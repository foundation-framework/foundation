package log

import (
	"go.uber.org/zap"
)

type Logger interface {
	Named(name string) Logger
	With(keysAndValues ...interface{}) Logger

	Debug(msg string, keysAndValues ...interface{})
	Info(msg string, keysAndValues ...interface{})
	Warn(msg string, keysAndValues ...interface{})
	Error(msg string, keysAndValues ...interface{})
	DPanic(msg string, keysAndValues ...interface{})
	Panic(msg string, keysAndValues ...interface{})
	Fatal(msg string, keysAndValues ...interface{})
}

type logger struct {
	logger *zap.SugaredLogger
}

func (d *logger) Named(name string) Logger {
	return &logger{logger: d.logger.Named(name)}
}

func (d *logger) With(keysAndValues ...interface{}) Logger {
	return &logger{logger: d.logger.With(keysAndValues)}
}

func (d *logger) Debug(msg string, keysAndValues ...interface{}) {
	d.logger.Debugw(msg, keysAndValues...)
}

func (d *logger) Info(msg string, keysAndValues ...interface{}) {
	d.logger.Infow(msg, keysAndValues...)
}

func (d *logger) Warn(msg string, keysAndValues ...interface{}) {
	d.logger.Warnw(msg, keysAndValues...)
}

func (d *logger) Error(msg string, keysAndValues ...interface{}) {
	d.logger.Errorw(msg, keysAndValues...)
}

func (d *logger) DPanic(msg string, keysAndValues ...interface{}) {
	d.logger.DPanicw(msg, keysAndValues...)
}

func (d *logger) Panic(msg string, keysAndValues ...interface{}) {
	d.logger.Panicw(msg, keysAndValues...)
}

func (d *logger) Fatal(msg string, keysAndValues ...interface{}) {
	d.logger.Fatalw(msg, keysAndValues...)
}
