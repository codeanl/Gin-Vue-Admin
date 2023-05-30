package router

import (
	"gin-vue-admin/gin-vue-admin/api"
	"gin-vue-admin/gin-vue-admin/middleware"
	"github.com/gin-gonic/gin"
)

func InitOperationLogRoutes(r *gin.RouterGroup) gin.IRoutes {
	router := r.Group("/log")
	// 开启jwt认证中间件
	router.Use(middleware.InitAuth())
	// 开启casbin鉴权中间件
	router.Use(middleware.CasbinMiddleware())
	{
		router.GET("/operation/list", api.GetOperationLogs)
		router.DELETE("/operation/delete/batch", api.BatchDeleteOperationLogByIds)
	}
	return r
}
