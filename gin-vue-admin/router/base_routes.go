package router

import (
	"gin-vue-admin/gin-vue-admin/api"
	"github.com/gin-gonic/gin"
)

// 注册基础路由
func InitBaseRoutes(r *gin.RouterGroup) gin.IRoutes {
	router := r.Group("/user")
	{
		//// 登录登出刷新token无需鉴权
		router.POST("/login", api.Login)
		router.POST("/logout", api.Logout)
	}
	return r
}
