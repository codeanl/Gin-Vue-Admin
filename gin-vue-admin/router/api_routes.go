package router

import (
	"gin-vue-admin/gin-vue-admin/api"
	"gin-vue-admin/gin-vue-admin/middleware"
	"github.com/gin-gonic/gin"
)

func InitApiRoutes(r *gin.RouterGroup) gin.IRoutes {
	router := r.Group("/api")
	// 开启jwt认证中间件
	router.Use(middleware.InitAuth())
	// 开启casbin鉴权中间件
	router.Use(middleware.CasbinMiddleware())
	{
		router.GET("/list", api.GetApis)
		router.GET("/tree", api.GetApiTree)
		router.POST("/create", api.CreateApi)
		router.POST("/update", api.UpdateApiById)
		router.DELETE("/delete/batch", api.BatchDeleteApiByIds)
	}
	return r
}
