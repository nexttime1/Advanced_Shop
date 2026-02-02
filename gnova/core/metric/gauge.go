package metric

import (
	prom "github.com/prometheus/client_golang/prometheus"
)

type (
	// GaugeVecOpts 是VectorOpts的别名，语义化（Gauge的配置）
	GaugeVecOpts VectorOpts

	// GaugeVec 接口：比Counter多了Set方法（直接设置值）
	GaugeVec interface {
		Set(v float64, labels ...string) // 直接设置值
		Inc(labels ...string)            // +1
		Add(v float64, labels ...string) // 加任意值（可正可负）
	}

	// promGaugeVec 私有实现体，持有官方库的GaugeVec
	promGaugeVec struct {
		gauge *prom.GaugeVec
	}
)

// NewGaugeVec returns a GaugeVec.
func NewGaugeVec(cfg *GaugeVecOpts) GaugeVec {
	if cfg == nil {
		return nil
	}

	vec := prom.NewGaugeVec(
		prom.GaugeOpts{
			Namespace: cfg.Namespace,
			Subsystem: cfg.Subsystem,
			Name:      cfg.Name,
			Help:      cfg.Help,
		}, cfg.Labels)
	prom.MustRegister(vec)
	gv := &promGaugeVec{
		gauge: vec,
	}

	return gv
}

func (gv *promGaugeVec) Inc(labels ...string) {
	gv.gauge.WithLabelValues(labels...).Inc()
}

func (gv *promGaugeVec) Add(v float64, labels ...string) {
	gv.gauge.WithLabelValues(labels...).Add(v)
}

func (gv *promGaugeVec) Set(v float64, labels ...string) {
	gv.gauge.WithLabelValues(labels...).Set(v)
}
