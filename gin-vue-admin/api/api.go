package api

import (
	"gin-vue-admin/gin-vue-admin/service"
	"gin-vue-admin/gin-vue-admin/util"
	"github.com/gin-gonic/gin"
)

// 获取接口列表
func GetApis(c *gin.Context) {
	var service service.ApiService
	if err := c.ShouldBind(&service); err == nil {
		res := service.GetApis(c)
		c.JSON(200, res)
	} else {
		c.JSON(400, err)
	}
}

// // 获取接口树
func GetApiTree(c *gin.Context) {
	var service service.ApiService
	if err := c.ShouldBind(&service); err == nil {
		res := service.GetApiTree(c)
		c.JSON(200, res)
	} else {
		c.JSON(400, err)
	}
}

// // 创建接口
func CreateApi(c *gin.Context) {
	var service service.ApiService
	claims, _ := util.ParseToken(c.GetHeader("Authorization"))
	if err := c.ShouldBind(&service); err == nil {
		res := service.CreateApi(c, claims.ID)
		c.JSON(200, res)
	} else {
		c.JSON(400, err)
	}
}

// UpdateApiById 更新接口
func UpdateApiById(c *gin.Context) {
	var service service.ApiService
	claims, _ := util.ParseToken(c.GetHeader("Authorization"))
	if err := c.ShouldBind(&service); err == nil {
		res := service.UpdateApiById(c, claims.ID)
		c.JSON(200, res)
	} else {
		c.JSON(400, err)
	}
}

// BatchDeleteApiByIds 批量删除接口
func BatchDeleteApiByIds(c *gin.Context) {
	var service service.ApiService
	claims, _ := util.ParseToken(c.GetHeader("Authorization"))
	if err := c.ShouldBind(&service); err == nil {
		res := service.BatchDeleteApiByIds(c, claims.ID)
		c.JSON(200, res)
	} else {
		c.JSON(400, err)
	}
}
