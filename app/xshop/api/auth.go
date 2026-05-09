package admin

import (
	"Advanced_Shop/app/pkg/options"
	"Advanced_Shop/gnova/server/restserver/middlewares"
	"Advanced_Shop/gnova/server/restserver/middlewares/auth"
	"github.com/gin-gonic/gin"
	"log"

	ginjwt "github.com/appleboy/gin-jwt/v2"
)

func newJWTAuth(opts *options.JwtOptions) middlewares.AuthStrategy {
	gjwt, err := ginjwt.New(&ginjwt.GinJWTMiddleware{
		Realm:            opts.Realm,
		SigningAlgorithm: "HS256",
		Key:              []byte(opts.Key),
		Timeout:          opts.Timeout,
		MaxRefresh:       opts.MaxRefresh,
		LogoutResponse:   func(c *gin.Context, code int) { c.JSON(code, nil) },
		IdentityHandler:  claimHandlerFunc,
		IdentityKey:      middlewares.KeyUserID,
		TokenLookup:      "header: Authorization:, query: token, cookie: jwt",
		TokenHeadName:    "Bearer",
		// 关闭默认未登录响应 改为自定义函数返回err
		Unauthorized: func(c *gin.Context, code int, message string) {
			log.Printf("JWT未授权：code=%d, message=%s", code, message) // 加日志
			c.JSON(code, gin.H{
				"code": code,
				"msg":  message,
			})
		},
	})
	if err != nil {
		panic("JWT中间件初始化失败：" + err.Error())
	}
	return auth.NewJWTStrategy(*gjwt)
}

func claimHandlerFunc(c *gin.Context) interface{} {
	claims := ginjwt.ExtractClaims(c)
	// 存入userid 和 role   JWT解析后数值为float64，暂存
	c.Set(middlewares.KeyUserID, claims[middlewares.KeyUserID])
	c.Set(middlewares.KeyRole, claims[middlewares.KeyRole])
	return claims[middlewares.KeyUserID]
}
