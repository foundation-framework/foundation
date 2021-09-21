package metrics

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type prometheusInstance struct {
	registry *prometheus.Registry
}

func NewPrometheusInstance() Instance {
	return &prometheusInstance{
		registry: prometheus.NewRegistry(),
	}
}

func (m *prometheusInstance) NewCounter(name, description string, labels ...string) Counter {
	vec := prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: name,
		Help: description,
	}, labels)

	return &prometheusCounter{vec: vec}
}

func (m *prometheusInstance) NewGauge(name, description string, labels ...string) Gauge {
	vec := prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: name,
		Help: description,
	}, labels)

	return &prometheusGauge{vec: vec}
}

func (m *prometheusInstance) Handler() http.Handler {
	return promhttp.HandlerFor(
		m.registry,
		promhttp.HandlerOpts{},
	)
}

type prometheusCounter struct {
	vec    *prometheus.CounterVec
	metric prometheus.Counter
}

func (c *prometheusCounter) With(labels ...string) (Counter, error) {
	metric, err := c.vec.GetMetricWith(stringSlicePairs(labels))
	if err != nil {
		return nil, err
	}

	return &prometheusCounter{
		vec:    c.vec,
		metric: metric,
	}, nil
}

func (c *prometheusCounter) WithValues(values ...string) (Counter, error) {
	metric, err := c.vec.GetMetricWithLabelValues(values...)
	if err != nil {
		return nil, err
	}

	return &prometheusCounter{
		vec:    c.vec,
		metric: metric,
	}, nil
}

func (c *prometheusCounter) Inc() {
	if c.metric != nil {
		c.metric.Inc()
	}
}

func (c *prometheusCounter) Add(delta float64) {
	if c.metric != nil {
		c.metric.Add(delta)
	}
}

type prometheusGauge struct {
	vec    *prometheus.GaugeVec
	metric prometheus.Gauge
}

func (c *prometheusGauge) With(labels ...string) (Gauge, error) {
	metric, err := c.vec.GetMetricWith(stringSlicePairs(labels))
	if err != nil {
		return nil, err
	}

	return &prometheusGauge{
		vec:    c.vec,
		metric: metric,
	}, nil
}

func (c *prometheusGauge) WithValues(values ...string) (Gauge, error) {
	metric, err := c.vec.GetMetricWithLabelValues(values...)
	if err != nil {
		return nil, err
	}

	return &prometheusGauge{
		vec:    c.vec,
		metric: metric,
	}, nil
}

func (c *prometheusGauge) Inc() {
	if c.metric != nil {
		c.metric.Inc()
	}
}

func (c *prometheusGauge) Dec() {
	if c.metric != nil {
		c.metric.Dec()
	}
}

func (c *prometheusGauge) Add(delta float64) {
	if c.metric != nil {
		c.metric.Add(delta)
	}
}

func (c *prometheusGauge) Sub(delta float64) {
	if c.metric != nil {
		c.metric.Sub(delta)
	}
}
