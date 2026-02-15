package common

import (
	"Advanced_Shop/gnova/server/restserver/middlewares"
	"Advanced_Shop/pkg/errors"
	"github.com/gin-gonic/gin"
)

func GetAuthUser(c *gin.Context) (int32, int, error) {
	// 尝试获取userid
	userIdVal, exists := c.Get(middlewares.KeyUserID)
	if !exists {
		return 0, 0, errors.New("未登录")
	}
	// 转换userid为uint  JWT解析后默认是float64
	userIdFloat, ok := userIdVal.(float64)
	if !ok {
		return 0, 0, errors.New("用户ID格式错误")
	}
	userID := uint(userIdFloat)

	// 尝试获取role
	roleVal, exists := c.Get(middlewares.KeyRole)
	if !exists {
		return 0, 0, errors.New("用户角色未配置")
	}
	// 转换role为int
	roleFloat, ok := roleVal.(float64)
	if !ok {
		return 0, 0, errors.New("用户角色格式错误")
	}
	role := int(roleFloat)

	// 校验角色合法性
	if role != 1 && role != 2 {
		return 0, 0, errors.New("角色不合法（仅支持1=管理员/2=普通用户）")
	}

	// 认证成功
	return int32(userID), role, nil
}
