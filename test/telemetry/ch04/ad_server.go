package main

import (
	"Advanced_Shop/test/telemetry/ch03/server/model"
	"Advanced_Shop/test/telemetry/ch04/common"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/baggage"
	"go.opentelemetry.io/otel/propagation"

	"go.opentelemetry.io/otel/sdk/trace"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
	"time"

	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	"gorm.io/plugin/opentelemetry/tracing"
)

var tp *trace.TracerProvider

func Server(c *gin.Context) {

	ctx := c.Request.Context()
	propagator := otel.GetTextMapPropagator()
	tr := tp.Tracer("basic_name")
	sctx := propagator.Extract(ctx, propagation.HeaderCarrier(c.Request.Header))
	spanCtx, span := tr.Start(sctx, "server")
	defer span.End()
	bag := baggage.FromContext(sctx)
	userID := bag.Member("user.id").Value()
	fmt.Println(userID)

	dsn := "root:root@tcp(127.0.0.1:3306)/user-srv?charset=utf8mb4&parseTime=True&loc=Local"

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
	})
	if err != nil {
		panic(err)
	}

	if err := db.Use(tracing.NewPlugin()); err != nil {
		panic(err)
	}

	if err := db.WithContext(spanCtx).Model(model.UserModels{}).Where("id = ?", 1).First(&model.UserModels{}).Error; err != nil {
		panic(err)
	}
	time.Sleep(500 * time.Millisecond)
	c.JSON(200, gin.H{})
}

func main() {
	_, tp = common.InitTracerProvider(common.Options{
		Name:     "basic_name",
		Endpoint: "http://192.168.163.132:14268/api/traces",
		Sampler:  1,
		Batcher:  "jaeger",
	})
	r := gin.Default()
	r.Use(otelgin.Middleware("gin_Middleware"))
	r.GET("/", func(c *gin.Context) {

	})
	r.GET("/server", Server)
	r.Run(":8090")
}
