package metrics

import "net/http"

// Instance describes basic metrics provider
type Instance interface {

	// NewCounter creates counter metric
	// See Counter for more information
	NewCounter(name, description string, labels ...string) Counter

	// NewGauge creates gauge metric
	// See Gauge for more information
	NewGauge(name, description string, labels ...string) Gauge

	// Handler returns http.Handler for metrics export
	Handler() http.Handler
}

// Counter is a metric that represents a single numerical
// value which we can only increase
type Counter interface {

	// With returns Counter with provided labels and values. If label
	// doesn't exist this method will create it
	With(labelsAndValues ...string) (Counter, error)

	// WithValues returns Counter with labels that have provided values.
	WithValues(values ...string) (Counter, error)

	// Inc increments counter
	Inc()

	// Add increments counter by provided delta
	Add(delta float64)
}

// Gauge is the same as the Counter, but we can also decrease it
type Gauge interface {

	// With returns Gauge with provided labels and values. If label
	// doesn't exist this method will create it
	With(labelsAndValues ...string) (Gauge, error)

	// WithValues returns Counter with labels that have provided values.
	WithValues(values ...string) (Gauge, error)

	// Inc increments gauge
	Inc()

	// Dec decrements gauge
	Dec()

	// Add increments gauge by provided delta
	Add(delta float64)

	// Sub decrements gauge by provided delta
	Sub(delta float64)
}
