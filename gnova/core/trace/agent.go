package trace

import (
	"Advanced_Shop/pkg/log"
	"context"
	"sync"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.21.0"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	kindJaeger = "jaeger"
	kindOTLP   = "otlp"
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

	err := startAgent(context.Background(), o)
	if err != nil {
		log.Errorf("InitAgent failed: %v", err)
		return
	}
	agents[o.Endpoint] = struct{}{}
}

func startAgent(ctx context.Context, o Options) error {
	// 1. 配置 Resource
	res, err := resource.New(ctx,
		resource.WithAttributes(
			semconv.ServiceNameKey.String(o.Name),
		),
	)
	if err != nil {
		return err
	}

	// 2. 创建 OTLP gRPC Exporter，连接到 otel-collector
	conn, err := grpc.NewClient(
		o.Endpoint, // 例如 "otel-collector:4317"
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return err
	}

	exp, err := otlptrace.New(ctx,
		otlptracegrpc.NewClient(
			otlptracegrpc.WithGRPCConn(conn),
		),
	)
	if err != nil {
		return err
	}

	// 3. 配置 TracerProvider
	tp := sdktrace.NewTracerProvider(
		sdktrace.WithSampler(
			sdktrace.ParentBased(sdktrace.TraceIDRatioBased(o.Sampler)),
		),
		sdktrace.WithResource(res),
		sdktrace.WithBatcher(exp),
	)

	// 4. 注册全局 TracerProvider
	otel.SetTracerProvider(tp)

	// 5. 设置分布式传播器
	otel.SetTextMapPropagator(
		propagation.NewCompositeTextMapPropagator(
			propagation.TraceContext{},
			propagation.Baggage{},
		),
	)

	// 6. 全局错误处理
	otel.SetErrorHandler(otel.ErrorHandlerFunc(func(err error) {
		log.Errorf("[otel] error: %v", err)
	}))

	return nil
}
