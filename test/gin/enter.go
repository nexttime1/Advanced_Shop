package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

func myMiddleWare(c *gin.Context) {
	fmt.Println("自己的中间件")
}
func myMiddleWare111(c *gin.Context) {
	fmt.Println("自己的中间件111")
}
func myMiddleWare1(c *gin.Context) {
	fmt.Println("自己的中间件1")
}
func myMiddleWare11(c *gin.Context) {
	fmt.Println("测试后给根 11")
}

func main() {
	// 创建一个 gin Engine，本质上是一个 http Handler
	mux := gin.Default()
	// 注册中间件
	mux.Use(myMiddleWare)
	// 注册一个 path 为 /ping 的处理函数
	mux.POST("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, "pone")
	})

	group := mux.Group("/api")
	mux.Use(myMiddleWare11)
	group.Use(myMiddleWare1)
	group.POST("/t1", func(c *gin.Context) {
		c.JSON(http.StatusOK, "pone")
	})
	mux.Use(myMiddleWare111)
	// 运行 http 服务
	if err := mux.Run(":8080"); err != nil {
		panic(err)
	}
}
