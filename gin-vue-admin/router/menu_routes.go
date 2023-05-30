package router

import (
	"gin-vue-admin/gin-vue-admin/api"
	"gin-vue-admin/gin-vue-admin/middleware"
	"github.com/gin-gonic/gin"
)

func InitMenuRoutes(r *gin.RouterGroup) gin.IRoutes {
	router := r.Group("/menu")
	// 开启jwt认证中间件
	router.Use(middleware.InitAuth())
	// 开启casbin鉴权中间件
	router.Use(middleware.CasbinMiddleware())
	{
		router.GET("/tree", api.GetMenuTree)
		router.GET("/list", api.GetMenus)
		router.POST("/create", api.CreateMenu)
		router.POST("/update", api.UpdateMenuById)
		router.DELETE("/delete/batch", api.BatchDeleteMenuByIds)
		router.GET("/access/list", api.GetUserMenusByUserId)
		router.GET("/access/tree", api.GetUserMenuTreeByUserId)
	}
	return r
}
