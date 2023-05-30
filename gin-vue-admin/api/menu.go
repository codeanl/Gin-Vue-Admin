package api

import (
	"gin-vue-admin/gin-vue-admin/service"
	"gin-vue-admin/gin-vue-admin/util"
	"github.com/gin-gonic/gin"
)

// 获取菜单列表
func GetMenus(c *gin.Context) {
	var service service.MenuService
	if err := c.ShouldBind(&service); err == nil {
		res := service.GetMenus(c)
		c.JSON(200, res)
	} else {
		c.JSON(400, err)
	}
}

// 获取菜单树
func GetMenuTree(c *gin.Context) {
	var service service.MenuService
	if err := c.ShouldBind(&service); err == nil {
		res := service.GetMenuTree(c)
		c.JSON(200, res)
	} else {
		c.JSON(400, err)
	}
}

// // 创建菜单
func CreateMenu(c *gin.Context) {
	var service service.MenuService
	claims, _ := util.ParseToken(c.GetHeader("Authorization"))
	if err := c.ShouldBind(&service); err == nil {
		res := service.CreateMenu(c, claims.ID)
		c.JSON(200, res)
	} else {
		c.JSON(400, err)
	}
}

// 更新菜单
func UpdateMenuById(c *gin.Context) {
	var service service.MenuService
	claims, _ := util.ParseToken(c.GetHeader("Authorization"))
	if err := c.ShouldBind(&service); err == nil {
		res := service.UpdateMenuById(c, claims.ID)
		c.JSON(200, res)
	} else {
		c.JSON(400, err)
	}
}

// 批量删除菜单
func BatchDeleteMenuByIds(c *gin.Context) {
	var service service.MenuService
	if err := c.ShouldBind(&service); err == nil {
		res := service.BatchDeleteMenuByIds(c)
		c.JSON(200, res)
	} else {
		c.JSON(400, err)
	}
}

// 根据用户ID获取用户的可访问菜单列表
func GetUserMenusByUserId(c *gin.Context) {
	var service service.MenuService
	if err := c.ShouldBind(&service); err == nil {
		res := service.GetUserMenusByUserId(c)
		c.JSON(200, res)
	} else {
		c.JSON(400, err)
	}
}

// 根据用户ID获取用户的可访问菜单树
func GetUserMenuTreeByUserId(c *gin.Context) {
	var service service.MenuService
	if err := c.ShouldBind(&service); err == nil {
		res := service.GetUserMenuTreeByUserId(c)
		c.JSON(200, res)
	} else {
		c.JSON(400, err)
	}
}
