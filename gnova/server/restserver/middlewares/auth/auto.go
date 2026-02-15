package auth

import (
	"Advanced_Shop/pkg/common/core"
	"strings"

	"Advanced_Shop/gnova/code"
	"Advanced_Shop/gnova/server/restserver/middlewares"

	"Advanced_Shop/pkg/errors"
	"github.com/gin-gonic/gin"
)

const authHeaderCount = 2

// AutoStrategy defines authentication strategy which can automatically choose between Basic and Bearer
// according `Authorization` header.
type AutoStrategy struct {
	basic BasicStrategy
	jwt   JWTStrategy
}

var _ middlewares.AuthStrategy = &AutoStrategy{}

// NewAutoStrategy create auto strategy with basic strategy and jwt strategy.
func NewAutoStrategy(basic BasicStrategy, jwt JWTStrategy) AutoStrategy {
	return AutoStrategy{
		basic: basic,
		jwt:   jwt,
	}
}

// AuthFunc defines auto strategy as the gin authentication middleware.
func (a AutoStrategy) AuthFunc() gin.HandlerFunc {
	return func(c *gin.Context) {
		operator := middlewares.AuthOperator{}
		authHeader := strings.SplitN(c.Request.Header.Get("Authorization"), " ", 2)

		if len(authHeader) != authHeaderCount {
			core.WriteErrResponse(
				c,
				errors.WithCode(code.ErrInvalidAuthHeader, "Authorization header format is wrong."),
				nil,
			)
			c.Abort()

			return
		}

		switch authHeader[0] {
		case "Basic":
			operator.SetStrategy(a.basic)
		case "Bearer":
			operator.SetStrategy(a.jwt)
			// a.JWT.MiddlewareFunc()(c)
		default:
			core.WriteErrResponse(c, errors.WithCode(code.ErrSignatureInvalid, "unrecognized Authorization header."), nil)
			c.Abort()

			return
		}

		operator.AuthFunc()(c)

		c.Next()
	}
}
