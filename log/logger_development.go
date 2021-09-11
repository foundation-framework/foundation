package log

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type developmentLogger struct {
	logger *zap.SugaredLogger
}

func NewDevelopmentLogger(output ...string) (Logger, error) {
	if len(output) == 0 {
		output = []string{"stdout"}
	}

	config := zap.Config{
		Level:            zap.NewAtomicLevelAt(zapcore.DebugLevel),
		Development:      true,
		Encoding:         "console",
		EncoderConfig:    zap.NewDevelopmentEncoderConfig(),
		OutputPaths:      output,
		ErrorOutputPaths: output,
	}

	logger, err := config.Build(zap.AddCallerSkip(1))
	if err != nil {
		return nil, err
	}

	return &developmentLogger{logger: logger.Sugar()}, nil
}

func (d *developmentLogger) Named(name string) Logger {
	return &developmentLogger{logger: d.logger.Named(name)}
}

func (d *developmentLogger) With(keysAndValues ...interface{}) Logger {
	return &developmentLogger{logger: d.logger.With(keysAndValues)}
}

func (d *developmentLogger) Debug(msg string, keysAndValues ...interface{}) {
	d.logger.Debugw(msg, keysAndValues...)
}

func (d *developmentLogger) Info(msg string, keysAndValues ...interface{}) {
	d.logger.Infow(msg, keysAndValues...)
}

func (d *developmentLogger) Warn(msg string, keysAndValues ...interface{}) {
	d.logger.Warnw(msg, keysAndValues...)
}

func (d *developmentLogger) Error(msg string, keysAndValues ...interface{}) {
	d.logger.Errorw(msg, keysAndValues...)
}

func (d *developmentLogger) DPanic(msg string, keysAndValues ...interface{}) {
	d.logger.DPanicw(msg, keysAndValues...)
}

func (d *developmentLogger) Panic(msg string, keysAndValues ...interface{}) {
	d.logger.Panicw(msg, keysAndValues...)
}

func (d *developmentLogger) Fatal(msg string, keysAndValues ...interface{}) {
	d.logger.Fatalw(msg, keysAndValues...)
}
