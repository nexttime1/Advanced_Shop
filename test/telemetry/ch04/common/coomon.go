package common

import (
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/exporters/zipkin"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.38.0"

	"Advanced_Shop/pkg/log"
)

const (
	kindJaeger = "jaeger"
	kindZipkin = "zipkin"
)

type Options struct {
	Name     string  `json:"name"`     // 服务名称，会显示在 Jaeger/Zipkin UI 中
	Endpoint string  `json:"endpoint"` // 收集器地址 如果是 jaeger 我们填的就是http://jaeger:14268/api/traces
	Sampler  float64 `json:"sampler"`  // 采样率，0.0~1.0
	Batcher  string  `json:"batcher"`  // 后端类型: "jaeger" 或 "zipkin" 这个字段决定 Endpoint字段填什么
}

func InitTracerProvider(o Options) (error, *trace.TracerProvider) {
	var sexp trace.SpanExporter
	var err error

	opts := []trace.TracerProviderOption{
		trace.WithSampler(trace.ParentBased(trace.TraceIDRatioBased(o.Sampler))),
		trace.WithResource(resource.NewSchemaless(semconv.ServiceNameKey.String(o.Name))),
	}

	if len(o.Endpoint) > 0 {
		switch o.Batcher {
		case kindJaeger:
			sexp, err = jaeger.New(jaeger.WithCollectorEndpoint(jaeger.WithEndpoint(o.Endpoint)))
			if err != nil {
				return err, nil
			}
		case kindZipkin:
			sexp, err = zipkin.New(o.Endpoint)
			if err != nil {
				return err, nil
			}
		}
		opts = append(opts, trace.WithBatcher(sexp))
	}

	tp := trace.NewTracerProvider(opts...)
	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))
	otel.SetErrorHandler(otel.ErrorHandlerFunc(func(err error) {
		log.Errorf("[otel] error: %v", err)
	}))
	return nil, tp
}
