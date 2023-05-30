package api

import (
	"gin-vue-admin/gin-vue-admin/service"
	"gin-vue-admin/gin-vue-admin/util"
	"github.com/gin-gonic/gin"
)

// 获取角色列表
func GetRoles(c *gin.Context) {
	var service service.RoleService
	if err := c.ShouldBind(&service); err == nil {
		res := service.GetRoles(c)
		c.JSON(200, res)
	} else {
		c.JSON(400, err)
	}
}

// 创建角色
func CreateRole(c *gin.Context) {
	var service service.RoleService
	chaim, _ := util.ParseToken(c.GetHeader("Authorization"))
	if err := c.ShouldBind(&service); err == nil {
		res := service.CreateRole(c, chaim.ID)
		c.JSON(200, res)
	} else {
		c.JSON(400, err)
	}
}

// 更新角色
func UpdateRoleById(c *gin.Context) {
	var service service.RoleService
	chaim, _ := util.ParseToken(c.GetHeader("Authorization"))
	if err := c.ShouldBind(&service); err == nil {
		res := service.UpdateRoleById(c, chaim.ID)
		c.JSON(200, res)
	} else {
		c.JSON(400, err)
	}
}

// 获取角色的权限菜单
func GetRoleMenusById(c *gin.Context) {
	var service service.RoleService
	if err := c.ShouldBind(&service); err == nil {
		res := service.GetRoleMenusById(c)
		c.JSON(200, res)
	} else {
		c.JSON(400, err)
	}
}

// 更新角色的权限菜单
func UpdateRoleMenusById(c *gin.Context) {
	var service service.RoleService
	chaim, _ := util.ParseToken(c.GetHeader("Authorization"))
	if err := c.ShouldBind(&service); err == nil {
		res := service.UpdateRoleMenusById(c, chaim.ID)
		c.JSON(200, res)
	} else {
		c.JSON(400, err)
	}
}

// 根据角色关键字获取角色的权限接口
func GetRoleApisById(c *gin.Context) {
	var service service.RoleService
	if err := c.ShouldBind(&service); err == nil {
		res := service.GetRoleApisById(c)
		c.JSON(200, res)
	} else {
		c.JSON(400, err)
	}
}

// 更新角色的权限接口
func UpdateRoleApisById(c *gin.Context) {
	var service service.RoleService
	chaim, _ := util.ParseToken(c.GetHeader("Authorization"))
	if err := c.ShouldBind(&service); err == nil {
		res := service.UpdateRoleApisById(c, chaim.ID)
		c.JSON(200, res)
	} else {
		c.JSON(400, err)
	}
}

// 批量删除角色
func BatchDeleteRoleByIds(c *gin.Context) {
	var service service.RoleService
	chaim, _ := util.ParseToken(c.GetHeader("Authorization"))
	if err := c.ShouldBind(&service); err == nil {
		res := service.BatchDeleteRoleByIds(c, chaim.ID)
		c.JSON(200, res)
	} else {
		c.JSON(400, err)
	}
}
