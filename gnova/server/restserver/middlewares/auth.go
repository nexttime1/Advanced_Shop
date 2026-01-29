package middlewares

import (
	"github.com/gin-gonic/gin"
)

type AuthStrategy interface {
	AuthFunc() gin.HandlerFunc
}

// AuthOperator 自主选择策略
type AuthOperator struct {
	strategy AuthStrategy
}

func (ao *AuthOperator) SetStrategy(strategy AuthStrategy) {
	ao.strategy = strategy
}

func (ao *AuthOperator) AuthFunc() gin.HandlerFunc {
	return ao.strategy.AuthFunc()
}
