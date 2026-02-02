package metric

import (
	prom "github.com/prometheus/client_golang/prometheus"
)

type (
	// CounterVecOpts 是VectorOpts的别名，语义化命名（明确是Counter的配置）
	CounterVecOpts VectorOpts

	// CounterVec 接口：定义Counter向量的操作方法，屏蔽底层实现
	CounterVec interface {
		// Inc 给指定标签的指标+1
		Inc(labels ...string)
		// Add 给指定标签的指标加任意浮点数（比如一次性加5）
		Add(v float64, labels ...string)
	}

	// promCounterVec 是CounterVec接口的实现体（私有结构体，外部不能直接访问）
	promCounterVec struct {
		counter *prom.CounterVec // 持有官方库的CounterVec实例
	}
)

// NewCounterVec 对外暴露的构造函数，返回CounterVec接口
func NewCounterVec(cfg *CounterVecOpts) CounterVec {

	if cfg == nil {
		return nil
	}

	//  调用官方库创建CounterVec实例
	vec := prom.NewCounterVec(prom.CounterOpts{
		// 把自定义配置映射到官方库的CounterOpts
		Namespace: cfg.Namespace,
		Subsystem: cfg.Subsystem,
		Name:      cfg.Name,
		Help:      cfg.Help,
	}, cfg.Labels) // 第二个参数是标签列表

	// 注册指标到Prometheus默认注册表（必须注册，否则Prometheus抓不到数据）
	prom.MustRegister(vec)

	// 4. 封装成自定义的promCounterVec，返回接口类型
	cv := &promCounterVec{
		counter: vec,
	}

	return cv
}
func (cv *promCounterVec) Inc(labels ...string) {
	cv.counter.WithLabelValues(labels...).Inc()
}

// Add  一次性加5（比如批量统计） reqCounter.Add(5, "GET", "404")
func (cv *promCounterVec) Add(v float64, labels ...string) {
	cv.counter.WithLabelValues(labels...).Add(v)
}
