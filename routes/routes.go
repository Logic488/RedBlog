package routes

import (
	"RedBlog/controller"
	"RedBlog/logger"
	"RedBlog/middleware"
	"net/http"

	"github.com/gin-gonic/gin"
)

func SetUp(mode string) *gin.Engine {

	if mode == gin.ReleaseMode {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.New()
	r.Use(logger.GinLogger(), logger.GinRecovery(true))

	v1 := r.Group("/api/v1")
	//注册业务路由
	v1.POST("signup", controller.SignUpHandler)
	//用户登录路由
	v1.POST("login", controller.LogInHandler)

	v1.Use(middleware.JWTAuthMiddleware()) //应用JWT认证中间件

	{
		//查询所有社区
		v1.GET("community", controller.CommunityHandler)
		//查询社区详情
		v1.GET("community/:id", controller.CommunityDetailHandler)

		//发表帖子
		v1.POST("post", controller.CreatePostHandler)
		//查询帖子详情
		v1.GET("post/:id", controller.PostDetailHandler)
		//查询帖子列表
		v1.GET("post/", controller.GetPostListHandler)
	}

	//r.GET("ping", middleware.JWTAuthMiddleware(), func(c *gin.Context) {
	//	var isLogin = true
	//	c.Request.Header.Get("Authorization")
	//	if isLogin {
	//		//如果是登录的用户，判断请求头中是否有有效的JWT Token
	//		c.String(http.StatusOK, "pong")
	//	} else {
	//		//否则就直接返回请登录
	//		c.String(http.StatusOK, "请登录")
	//	}
	//})

	r.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "ok")
	})

	r.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"code": 404,
		})
	})

	return r
}
