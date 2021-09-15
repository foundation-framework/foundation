package log

import (
	"fmt"

	"go.uber.org/zap"
)

func init() {
	if err := zap.RegisterSink(udpScheme, newUdpSink); err != nil {
		panic(fmt.Errorf("unexpected error: %s", err.Error()))
	}
}
