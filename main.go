package main

import (
	. "gorm-test/user"

	"github.com/gin-gonic/gin"
)

func main() {

	router := gin.Default()

	// 注册和登录处理函数需要传递db参数
	router.POST("/register", func(c *gin.Context) {
		Register(c)
	})

	router.POST("/login", func(c *gin.Context) {
		Login(c)
	})

	//
	protectGroup := router.Group("user").Use(JwtMiddleware())
	protectGroup.POST("/createpost", CreatePost)
	protectGroup.GET("/listposts", ListPosts)
	protectGroup.POST("/updatepost", UpdatePost)

	router.Run()

}
