package main

import (
	"Advanced_Shop/pkg/log"
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.12.0"

	"github.com/valyala/fasthttp"
)

const (
	traceName = "xshop-otel"
)

var tp *trace.TracerProvider

func tracerProvider() error {
	url := "http://127.0.0.1:14268/api/traces"
	jexp, err := jaeger.New(jaeger.WithCollectorEndpoint(jaeger.WithEndpoint(url)))
	if err != nil {
		panic(err)
	}

	tp = trace.NewTracerProvider(
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
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))
	return nil
}

func funcA(ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()
	tr := otel.Tracer("traceName")
	spanCtx, span := tr.Start(ctx, "func-a")
	span.SetAttributes(attribute.String("name", "funA"))

	type _LogStruct struct {
		CurrentTime time.Time `json:"current_time"`
		PassWho     string    `json:"pass_who"`
		Name        string    `json:"name"`
	}

	logTest := _LogStruct{
		CurrentTime: time.Now(),
		PassWho:     "bobby",
		Name:        "func-a",
	}

	b, _ := json.Marshal(logTest)
	//日志会发送到jaeger
	log.InfofC(spanCtx, "this is funca log: %s", string(b))

	span.SetAttributes(attribute.Key("这是测试日志的key").String(string(b)))
	time.Sleep(time.Second)
	span.End()
}

func funcB(ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()
	tr := otel.Tracer("traceName")
	spanCtx, span := tr.Start(ctx, "func-b")
	fmt.Println("trace:", span.SpanContext().TraceID(), span.SpanContext().SpanID())
	time.Sleep(time.Second)

	req := fasthttp.AcquireRequest()
	req.SetRequestURI("http://127.0.0.1:8090/server")
	req.Header.SetMethod("GET")

	//拿起传播器
	p := otel.GetTextMapPropagator()
	//包裹
	headers := make(map[string]string)
	p.Inject(spanCtx, propagation.MapCarrier(headers))

	for key, value := range headers {
		req.Header.Set(key, value)
	}
	//req.Header.Set("trace-id", span.SpanContext().TraceID().String())
	//req.Header.Set("span-id", span.SpanContext().SpanID().String())

	fclient := fasthttp.Client{}
	fres := fasthttp.Response{}
	_ = fclient.Do(req, &fres)

	log.InfofC(spanCtx, "this is funcB log: %s", "imooc")

	span.End()
}

func main() {
	_ = tracerProvider()
	ctx, cancel := context.WithCancel(context.Background())
	defer func(ctx context.Context) {
		ctx, cancel = context.WithTimeout(ctx, time.Second*5)
		defer cancel()
		if err := tp.Shutdown(ctx); err != nil {
			panic(err)
		}
	}(ctx)

	tr := otel.Tracer(traceName)
	spanCtx, span := tr.Start(ctx, "func-main")

	wg := &sync.WaitGroup{}
	wg.Add(2)
	go funcA(spanCtx, wg)
	go funcB(spanCtx, wg)

	span.AddEvent("this is an event")
	time.Sleep(time.Second)
	wg.Wait()
	span.End()
}
