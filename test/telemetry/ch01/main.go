package main

import (
	"context"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.38.0"
)

func main() {
	url := "http://127.0.0.1:14268/api/traces"
	jexp, err := jaeger.New(jaeger.WithCollectorEndpoint(jaeger.WithEndpoint(url)))
	if err != nil {
		panic(err)
	}

	tp := trace.NewTracerProvider(
		trace.WithBatcher(jexp),
		trace.WithResource(
			resource.NewWithAttributes(
				semconv.SchemaURL,
				semconv.ServiceNameKey.String("xshop-user"),
				attribute.String("environment", "dev"),
				attribute.Int("ID", 1),
			),
		),
	)
	otel.SetTracerProvider(tp)

	ctx, cancel := context.WithCancel(context.Background())
	defer func(ctx context.Context) {
		ctx, cancel = context.WithTimeout(ctx, time.Second*5)
		defer cancel()
		if err := tp.Shutdown(ctx); err != nil {
			panic(err)
		}
	}(ctx)

	tr := otel.Tracer("xshop-otel")
	_, span := tr.Start(ctx, "func-main")
	var attrs []attribute.KeyValue
	attrs = append(attrs, attribute.String("key1", "value1"))
	attrs = append(attrs, attribute.Bool("key2", false))
	attrs = append(attrs, attribute.Int("key3", 123))
	attrs = append(attrs, attribute.StringSlice("key4", []string{"value4-1", "value4-2"}))

	span.SetAttributes(attrs...)

	span.AddEvent("this is an event")
	time.Sleep(time.Second)
	span.End()
}
