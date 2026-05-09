package middlewares

import (
	"Advanced_Shop/pkg/log"
	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
)

func TracingHandler(service string) gin.HandlerFunc {
	log.Infof("链路 service: %v", service)
	return otelgin.Middleware(service)
}
