package main

import (
	"GoStart/telemetry/ch03/server/model"
	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"log"
	"os"
	"time"

	semconv "go.opentelemetry.io/otel/semconv/v1.12.0"

	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	"gorm.io/plugin/opentelemetry/tracing"
)

var tp *trace.TracerProvider

const (
	traceName = "mxshop-otel"
)

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
				semconv.ServiceNameKey.String("mxshop-user"),
				attribute.String("environment", "dev"),
				attribute.Int("ID", 1),
			),
		),
	)
	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))
	return nil
}

func Server(c *gin.Context) {
	dsn := "root:root@tcp(127.0.0.1:3306)/mxshop_user_srv?charset=utf8mb4&parseTime=True&loc=Local"

	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			LogLevel: logger.Info, // Log level
			Colorful: true,        // Disable color
		},
	)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
		Logger: newLogger,
	})
	if err != nil {
		panic(err)
	}

	if err := db.Use(tracing.NewPlugin()); err != nil {
		panic(err)
	}
	//负责span的抽取和生成
	//ctx := c.Request.Context()
	//p := otel.GetTextMapPropagator()
	//tr := tp.Tracer(traceName)
	//sctx := p.Extract(ctx, propagation.HeaderCarrier(c.Request.Header))
	//spanCtx, span := tr.Start(sctx, "server")

	if err := db.WithContext(c.Request.Context()).Model(model.User{}).Where("id = ?", 1).First(&model.User{}).Error; err != nil {
		panic(err)
	}
	time.Sleep(500 * time.Millisecond)
	c.JSON(200, gin.H{})
}

//func Server(c *gin.Context) {
//
//	////负责span的抽取和生成
//	//us, _ := c.Value("trace").(string)
//	//if us == "" {
//	//	panic(us)
//	//}
//	//ctx := c.Request.Context()
//	////p := otel.GetTextMapPropagator()
//	////tr := tp.Tracer(traceName)
//	////sctx := p.Extract(ctx, propagation.HeaderCarrier(c.Request.Header))
//	//traceID := c.Request.Header.Get("trace-id")
//	//spanID := c.Request.Header.Get("span-id")
//	//
//	//tracid, _ := otelTrace.TraceIDFromHex(traceID)
//	//spanid, _ := otelTrace.SpanIDFromHex(spanID)
//	tr := otel.Tracer("traceName")
//	//spanCtx := otelTrace.NewSpanContext(otelTrace.SpanContextConfig{
//	//	TraceID:    tracid,
//	//	SpanID:     spanid,
//	//	TraceFlags: otelTrace.FlagsSampled, //这个如果不设置的话，是不会保存的
//	//	Remote:     true,
//	//})
//	//
//	//carrier := propagation.HeaderCarrier{}
//	//carrier.Set("trace-id", traceID)
//	//propagator := otel.GetTextMapPropagator()
//	//pctx := propagator.Extract(ctx, carrier)
//	//
//	//sct := otelTrace.ContextWithRemoteSpanContext(pctx, spanCtx)
//
//	_, span := tr.Start(sct, "server")
//	time.Sleep(600 * time.Millisecond)
//	span.End()
//	c.JSON(200, gin.H{})
//}

func main() {
	_ = tracerProvider()
	r := gin.Default()
	r.Use(otelgin.Middleware("my-server"))
	r.GET("/", func(c *gin.Context) {

	})
	r.GET("/server", Server)
	r.Run(":8090")
}
