package router

import (
	"gin-vue-admin/gin-vue-admin/api"
	"gin-vue-admin/gin-vue-admin/middleware"
	"github.com/gin-gonic/gin"
)

// 注册用户路由
func InitUserRoutes(r *gin.RouterGroup) gin.IRoutes {
	router := r.Group("/user")
	// 开启jwt认证中间件
	router.Use(middleware.InitAuth())
	// 开启casbin鉴权中间件
	router.Use(middleware.CasbinMiddleware())
	{
		router.GET("/info", api.GetMyInfo)
		//router.GET("/info", api.GetUserInfo)
		router.GET("/list", api.GetUsers)
		router.PUT("/changePwd", api.ChangePwd)
		router.POST("/create", api.CreateUser)
		router.POST("/update", api.UpdateUserById)
		router.DELETE("/delete/batch", api.BatchDeleteUserByIds)
	}
	return r
}
