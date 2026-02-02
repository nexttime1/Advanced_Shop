package main

import (
	"Advanced_Shop/test/telemetry/ch04/common"
	"context"
	"go.opentelemetry.io/otel/baggage"
	"golang.org/x/sync/errgroup"
	"time"

	"github.com/valyala/fasthttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/propagation"
)

func funcC(ctx context.Context) {
	tr := otel.Tracer("basic_name")
	_, span := tr.Start(ctx, "func-C")
	span.SetAttributes(attribute.String("FuncC的测试key", "内部key，不可传递"))
	time.Sleep(time.Second)
	span.End()
}

func funcD(ctx context.Context) {
	tr := otel.Tracer("basic_name")
	spanCtx, span := tr.Start(ctx, "func-D")
	userID, _ := baggage.NewMember("user.id", "10001")
	userRole, _ := baggage.NewMember("user.role", "1")
	bag, _ := baggage.New(userID, userRole)
	Newctx := baggage.ContextWithBaggage(spanCtx, bag)
	time.Sleep(time.Second)
	req := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(req)
	req.SetRequestURI("http://127.0.0.1:8090/server")
	req.Header.SetMethod("GET")
	//拿起传播器
	p := otel.GetTextMapPropagator()
	//包裹
	headers := make(map[string]string)
	p.Inject(Newctx, propagation.MapCarrier(headers))

	for key, value := range headers {
		req.Header.Set(key, value)
	}

	fclient := fasthttp.Client{}
	fres := fasthttp.Response{}
	// 隐藏错误
	_ = fclient.Do(req, &fres)

	span.End()
}

func main() {

	_, tp := common.InitTracerProvider(common.Options{
		Name:     "basic_name",
		Endpoint: "http://192.168.163.132:14268/api/traces",
		Sampler:  1,
		Batcher:  "jaeger",
	})

	ctx, cancel := context.WithCancel(context.Background())
	defer func(ctx context.Context) {
		ctx, cancel = context.WithTimeout(ctx, time.Second*5)
		defer cancel()
		if err := tp.Shutdown(ctx); err != nil {
			panic(err)
		}
	}(ctx)

	tr := otel.Tracer("basic_name")
	spanCtx, span := tr.Start(ctx, "func-main")

	gw, gctx := errgroup.WithContext(spanCtx)
	gw.Go(func() error {
		funcC(gctx)
		return nil
	})
	gw.Go(func() error {
		funcD(gctx)
		return nil
	})

	_ = gw.Wait()
	span.End()
}
