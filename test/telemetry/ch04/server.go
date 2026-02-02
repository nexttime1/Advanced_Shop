package main

import (
	"Advanced_Shop/test/telemetry/ch04/common"
	"context"
	"fmt"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"golang.org/x/sync/errgroup"
	"time"
)

func funcA(ctx context.Context) {
	tr := otel.Tracer("basic_name")
	_, span := tr.Start(ctx, "func-a")
	span.SetAttributes(attribute.String("测试Key", "独一无二"))
	span.AddEvent("完成 funcA")
	time.Sleep(time.Second)
	span.End()
}

func funcB(ctx context.Context) {
	tr := otel.Tracer("basic_name")
	_, span := tr.Start(ctx, "func-b")
	fmt.Println("trace:", span.SpanContext().TraceID(), span.SpanContext().SpanID())
	time.Sleep(time.Second)
	span.End()
}

func main() {
	_, tp := common.InitTracerProvider(common.Options{
		Name:     "basic_name",
		Endpoint: "http://192.168.163.132:14268/api/traces",
		Sampler:  1,
		Batcher:  "jaeger",
	})
	ctx := context.Background()
	defer func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_ = tp.Shutdown(ctx)
	}()

	tr := otel.Tracer("basic_name")
	spanCtx, span := tr.Start(ctx, "func-main")

	gw, gctx := errgroup.WithContext(spanCtx)
	gw.Go(func() error {
		funcA(gctx)
		return nil
	})
	gw.Go(func() error {
		funcB(gctx)
		return nil
	})

	_ = gw.Wait()
	span.End()

}
