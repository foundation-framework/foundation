package zaputil

import (
	"go.uber.org/zap"
)

func init() {
	if err := zap.RegisterSink(udpScheme, newUdpSink); err != nil {
		panic("zaputil: unexpected error: " + err.Error())
	}

	if err := zap.RegisterEncoder("logfmt", newLogfmtEncoder); err != nil {
		panic("zaputil: unexpected error: " + err.Error())
	}
}
