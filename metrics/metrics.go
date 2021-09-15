package metrics

import "github.com/intale-llc/foundation/transport/http"

type Metrics interface {
	NewCounter(name, description string, labels ...string) Counter
	NewGauge(name, description string, labels ...string) Gauge

	Handle(path string, listener http.Listener)
}

type Counter interface {
	With(labelsAndValues ...string) Counter
	WithValues(values ...string) Counter

	Inc()
	Add(delta float64)
}

type Gauge interface {
	With(labelsAndValues ...string) Gauge
	WithValues(values ...string) Gauge

	Inc()
	Dec()
	Add(delta float64)
	Sub(delta float64)
}
