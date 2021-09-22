package zaputil

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func NewLogfmtLogger(development bool, output ...string) (*zap.Logger, error) {
	if len(output) == 0 {
		output = []string{"stdout"}
	}

	config := zap.Config{
		Level:         zap.NewAtomicLevelAt(zapcore.DebugLevel),
		Encoding:      "logfmt",
		Development:   development,
		DisableCaller: true,
		EncoderConfig: zapcore.EncoderConfig{
			TimeKey:        "ts",
			LevelKey:       "level",
			NameKey:        "logger",
			CallerKey:      "caller",
			FunctionKey:    zapcore.OmitKey,
			MessageKey:     "msg",
			StacktraceKey:  "stack",
			LineEnding:     zapcore.DefaultLineEnding,
			EncodeLevel:    zapcore.LowercaseLevelEncoder,
			EncodeTime:     zapcore.RFC3339NanoTimeEncoder,
			EncodeDuration: zapcore.StringDurationEncoder,
			EncodeCaller:   zapcore.ShortCallerEncoder,
		},
		OutputPaths:      output,
		ErrorOutputPaths: output,
	}

	logger, err := config.Build()
	if err != nil {
		return nil, err
	}

	return logger, nil
}

func NewConsoleLogger(development bool, output ...string) (*zap.Logger, error) {
	if len(output) == 0 {
		output = []string{"stdout"}
	}

	config := zap.Config{
		Level:         zap.NewAtomicLevelAt(zapcore.InfoLevel),
		Encoding:      "console",
		Development:   development,
		DisableCaller: true,
		EncoderConfig: zapcore.EncoderConfig{
			TimeKey:        "T",
			LevelKey:       "L",
			NameKey:        "N",
			CallerKey:      "C",
			FunctionKey:    zapcore.OmitKey,
			MessageKey:     "M",
			StacktraceKey:  "S",
			LineEnding:     zapcore.DefaultLineEnding,
			EncodeLevel:    zapcore.CapitalLevelEncoder,
			EncodeTime:     zapcore.RFC3339NanoTimeEncoder,
			EncodeDuration: zapcore.SecondsDurationEncoder,
			EncodeCaller:   zapcore.ShortCallerEncoder,
		},
		OutputPaths:      output,
		ErrorOutputPaths: output,
	}

	logger, err := config.Build()
	if err != nil {
		return nil, err
	}

	return logger, nil
}

func NewJSONLogger(development bool, output ...string) (*zap.Logger, error) {
	if len(output) == 0 {
		output = []string{"stdout"}
	}

	config := zap.Config{
		Level:         zap.NewAtomicLevelAt(zapcore.InfoLevel),
		Encoding:      "json",
		Development:   development,
		DisableCaller: true,
		EncoderConfig: zapcore.EncoderConfig{
			TimeKey:        "ts",
			LevelKey:       "level",
			NameKey:        "logger",
			CallerKey:      "caller",
			FunctionKey:    zapcore.OmitKey,
			MessageKey:     "msg",
			StacktraceKey:  "stack",
			LineEnding:     zapcore.DefaultLineEnding,
			EncodeLevel:    zapcore.LowercaseLevelEncoder,
			EncodeTime:     zapcore.RFC3339NanoTimeEncoder,
			EncodeDuration: zapcore.SecondsDurationEncoder,
			EncodeCaller:   zapcore.ShortCallerEncoder,
		},
		OutputPaths:      output,
		ErrorOutputPaths: output,
	}

	logger, err := config.Build()
	if err != nil {
		return nil, err
	}

	return logger, nil
}
