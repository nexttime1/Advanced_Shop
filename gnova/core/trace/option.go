package trace

const TraceName = "mxshop"

type Options struct {
    Name     string  `json:"name"`      // 服务名称，会显示在 Jaeger/Zipkin UI 中
    Endpoint string  `json:"endpoint"`  // 收集器地址
    Sampler  float64 `json:"sampler"`   // 采样率，0.0~1.0
    Batcher  string  `json:"batcher"`   // 后端类型: "jaeger" 或 "zipkin"
}