package router

import (
	"gin-vue-admin/gin-vue-admin/api"
	"gin-vue-admin/gin-vue-admin/middleware"
	"github.com/gin-gonic/gin"
)

func InitRoleRoutes(r *gin.RouterGroup) gin.IRoutes {
	router := r.Group("/role")
	// 开启jwt认证中间件
	router.Use(middleware.InitAuth())
	// 开启casbin鉴权中间件
	router.Use(middleware.CasbinMiddleware())
	{
		router.GET("/list", api.GetRoles)
		router.POST("/create", api.CreateRole)
		router.POST("/update", api.UpdateRoleById)
		router.GET("/menus/get", api.GetRoleMenusById)
		router.POST("/menus/update", api.UpdateRoleMenusById)
		router.GET("/apis/get", api.GetRoleApisById)
		router.POST("/apis/update", api.UpdateRoleApisById)
		router.DELETE("/delete/batch", api.BatchDeleteRoleByIds)
	}
	return r
}
