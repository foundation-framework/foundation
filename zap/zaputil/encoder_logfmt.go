package zaputil

import (
	"go.uber.org/zap/zapcore"

	zaplogfmt "github.com/sykesm/zap-logfmt"
)

func newLogfmtEncoder(config zapcore.EncoderConfig) (zapcore.Encoder, error) {
	return zaplogfmt.NewEncoder(config), nil
}
