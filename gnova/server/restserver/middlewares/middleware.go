package middlewares

import (
	"github.com/gin-gonic/gin"
)

// Middlewares  全局默认
var Middlewares = defaultMiddlewares()

func defaultMiddlewares() map[string]gin.HandlerFunc {
	return map[string]gin.HandlerFunc{
		"recovery": gin.Recovery(),
		"cors":     Cors(),
		"context":  Context(),
	}
}
