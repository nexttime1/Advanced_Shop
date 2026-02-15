package middlewares

import "github.com/gin-gonic/gin"

const (
	UsernameKey = "username"
	KeyUserID   = "userid"
	UserIP      = "ip"
	KeyNickName = "nickname"
	KeyRole     = "role" // 角色（1=管理员，2=普通用户）
)

// Context 为每个请求添加上下文
func Context() gin.HandlerFunc {
	return func(c *gin.Context) {
		//TODO 扩展
		c.Next()
	}
}
