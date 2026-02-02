package metric

// VectorOpts 是所有指标向量（Vec）的基础配置结构体
type VectorOpts struct {
	Namespace string   // 命名空间（比如业务名：user_service、order_service）
	Subsystem string   // 子系统（比如模块名：api、db、cache）
	Name      string   // 指标名（比如：request_total、response_time）
	Help      string   // 指标说明（Prometheus UI里显示的注释，必须有）
	Labels    []string // 指标的标签（维度，比如：method、status、endpoint）
}
