package trace

import (
	"context"
	"sync"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/exporters/otlptrace/otlptracegrpc" // 新增
	"go.opentelemetry.io/otel/exporters/zipkin"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"Advanced_Shop/pkg/log"
)

const (
	kindJaeger = "jaeger"
	kindZipkin = "zipkin"
	kindOTLP   = "otlp" // 新增：推荐 K8s 下使用此类型
)

var (
	agents = make(map[string]struct{})
	lock   sync.Mutex
)

func InitAgent(o Options) {
	lock.Lock()
	defer lock.Unlock()

	_, ok := agents[o.Endpoint]
	if ok {
		return
	}
	err := startAgent(o)
	if err != nil {
		log.Errorf("InitAgent failed: %v", err)
		return
	}
	agents[o.Endpoint] = struct{}{}
}

func startAgent(o Options) error {
	var sexp trace.SpanExporter
	var err error
	ctx := context.Background()

	// 1. 设置 Resource 信息
	res, err := resource.New(ctx,
		resource.WithAttributes(
			semconv.ServiceNameKey.String(o.Name),
		),
	)
	if err != nil {
		return err
	}

	// 2. 根据 Batcher 类型初始化导出器
	if len(o.Endpoint) > 0 {
		switch o.Batcher {
		case kindOTLP:
			// 连接到 otel-collector (gRPC)
			sexp, err = otlptracegrpc.New(ctx,
				otlptracegrpc.WithInsecure(),
				otlptracegrpc.WithEndpoint(o.Endpoint),
				otlptracegrpc.WithDialOption(grpc.WithBlock()), // 可选：等待连接成功
				otlptracegrpc.WithTimeout(time.Second*5),
			)
		case kindJaeger:
			// 旧的 Jaeger 直接推送方式
			sexp, err = jaeger.New(jaeger.WithCollectorEndpoint(jaeger.WithEndpoint(o.Endpoint)))
		case kindZipkin:
			sexp, err = zipkin.New(o.Endpoint)
		}

		if err != nil {
			return err
		}
	}

	// 3. 配置 TracerProvider 选项
	opts := []trace.TracerProviderOption{
		trace.WithSampler(trace.ParentBased(trace.TraceIDRatioBased(o.Sampler))),
		trace.WithResource(res),
	}

	if sexp != nil {
		opts = append(opts, trace.WithBatcher(sexp))
	}

	// 4. 全局注册
	tp := trace.NewTracerProvider(opts...)
	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))

	otel.SetErrorHandler(otel.ErrorHandlerFunc(func(err error) {
		log.Errorf("[otel] error: %v", err)
	}))

	return nil
}
