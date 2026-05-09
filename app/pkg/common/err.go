package common

import (
	"Advanced_Shop/gnova/code"
	"Advanced_Shop/pkg/errors"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc/status"
)

type HandlerFunc func(c *gin.Context) error

func Wrapper(handler HandlerFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		err := handler(c)
		if err != nil {
			handleError(c, err)
		}
	}
}

// handleError 统一错误处理
func handleError(c *gin.Context, err error) {
	// 先检查是否是 web 层直接产生的 withCode 错误
	if errors.IsWithCode(err) {
		// 直接是 withCode，不需要转换（web层参数校验等错误）
		WriteErrResponse(c, err)
		return
	}

	// 检查是否是 gRPC 错误（从微服务返回的）
	if _, ok := status.FromError(err); ok {
		// 是 gRPC 错误，需要转换回 withCode
		customErr := errors.FromGrpcError(err)
		WriteErrResponse(c, customErr)
		return
	}

	// 其他未知错误 - 保底处理
	c.JSON(500, ErrResponse{
		Code:    code.ErrUnknown,
		Message: "internal server error",
		Detail:  err.Error(),
	})
}
func ErrorRecovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				var customErr error

				// 判断 panic 的类型
				switch e := err.(type) {
				case error:
					customErr = errors.FromGrpcError(e)
				case string:
					customErr = errors.WithCode(100002, e)
				default:
					customErr = errors.WithCode(100002, "unknown panic")
				}

				// 获取错误码信息
				coder := errors.ParseCoder(customErr)

				// 构造响应
				c.JSON(coder.HTTPStatus(), gin.H{
					"code":      coder.Code(),
					"message":   coder.String(),
					"reference": coder.Reference(),
				})

				c.Abort()
			}
		}()

		c.Next()
	}
}
