package metric

import (
	prom "github.com/prometheus/client_golang/prometheus"
)

type (
	// HistogramVecOpts 扩展了VectorOpts，多了Buckets（桶配置）
	HistogramVecOpts struct {
		Namespace string    // 基础字段（复用）
		Subsystem string    // 基础字段
		Name      string    // 基础字段
		Help      string    // 基础字段
		Labels    []string  // 基础字段
		Buckets   []float64 // 直方图的桶（核心扩展字段）
	}

	// HistogramVec 接口：只有Observe方法（观测一个数值）
	HistogramVec interface {
		Observe(v int64, labels ...string) // 观测值（这里封装成int64，官方是float64）
	}

	// promHistogramVec 私有实现体，持有官方库的HistogramVec
	promHistogramVec struct {
		histogram *prom.HistogramVec
	}
)

func NewHistogramVec(cfg *HistogramVecOpts) HistogramVec {
	if cfg == nil {
		return nil
	}

	// 映射到官方库的HistogramOpts（多了Buckets）
	vec := prom.NewHistogramVec(prom.HistogramOpts{
		Namespace: cfg.Namespace,
		Subsystem: cfg.Subsystem,
		Name:      cfg.Name,
		Help:      cfg.Help,
		Buckets:   cfg.Buckets, // 桶配置是Histogram的核心
	}, cfg.Labels)
	prom.MustRegister(vec)
	hv := &promHistogramVec{
		histogram: vec,
	}

	return hv
}

// Observe 观测一个数值（v是int64，内部转成float64）
func (hv *promHistogramVec) Observe(v int64, labels ...string) {
	hv.histogram.WithLabelValues(labels...).Observe(float64(v))
}
