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

// 作用：gin-jwt 解析完 Token 后，自动执行这个函数
func claimHandlerFunc(c *gin.Context) interface{} {
	// 1. 从 gin 上下文提取 JWT 载荷（map[string]interface{}）
	claims := ginjwt.ExtractClaims(c)
	// 2. 把 userid、role 存到 Gin 上下文，给业务接口用
	c.Set(middlewares.KeyUserID, claims[middlewares.KeyUserID])
	c.Set(middlewares.KeyRole, claims[middlewares.KeyRole])
	// 3. 返回用户ID（给 gin-jwt 内部使用）
	return claims[middlewares.KeyUserID]
}
