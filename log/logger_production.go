package log

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func NewProductionLogger(output ...string) (Logger, error) {
	if len(output) == 0 {
		output = []string{"stdout"}
	}

	config := zap.Config{
		Level:    zap.NewAtomicLevelAt(zapcore.InfoLevel),
		Encoding: "json",
		EncoderConfig: zapcore.EncoderConfig{
			TimeKey:        "timestamp",
			LevelKey:       "level",
			NameKey:        "logger",
			CallerKey:      "caller",
			FunctionKey:    zapcore.OmitKey,
			MessageKey:     "msg",
			StacktraceKey:  "stacktrace",
			LineEnding:     zapcore.DefaultLineEnding,
			EncodeLevel:    zapcore.LowercaseLevelEncoder,
			EncodeTime:     zapcore.ISO8601TimeEncoder,
			EncodeDuration: zapcore.SecondsDurationEncoder,
			EncodeCaller:   zapcore.ShortCallerEncoder,
		},
		OutputPaths:      output,
		ErrorOutputPaths: output,
	}

	lg, err := config.Build(zap.AddCallerSkip(1))
	if err != nil {
		return nil, err
	}

	return &logger{logger: lg.Sugar()}, nil
}
