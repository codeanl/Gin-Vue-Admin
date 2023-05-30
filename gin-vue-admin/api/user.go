package api

import (
	"gin-vue-admin/gin-vue-admin/service"
	"gin-vue-admin/gin-vue-admin/util"
	"github.com/gin-gonic/gin"
)

// 用户登录
func Login(c *gin.Context) {
	var service service.UserService
	if err := c.ShouldBind(&service); err == nil {
		res := service.Login(c)
		c.JSON(200, res)
	} else {
		c.JSON(400, err)
	}
}

// 退出登录
func Logout(c *gin.Context) {
	var service service.UserService
	Authorization := c.GetHeader("Authorization")
	if err := c.ShouldBind(&service); err == nil {
		res := service.Logout(c, Authorization)
		c.JSON(200, res)
	} else {
		c.JSON(400, err)
	}
}

// 获取自己的详细信息
func GetMyInfo(c *gin.Context) {
	var service service.UserService
	claims, _ := util.ParseToken(c.GetHeader("Authorization"))
	if err := c.ShouldBind(&service); err == nil {
		res := service.GetMyInfo(c, claims.ID)
		c.JSON(200, res)
	} else {
		c.JSON(400, err)
	}
}

// 获取用户列表
func GetUsers(c *gin.Context) {
	var service service.UserService
	if err := c.ShouldBind(&service); err == nil {
		res := service.GetUsers(c)
		c.JSON(200, res)
	} else {
		c.JSON(400, err)
	}
}

// 更新用户登录密码
func ChangePwd(c *gin.Context) {
	oldpassword := c.PostForm("oldpassword")
	newpassword := c.PostForm("newpassword")
	var service service.UserService
	claims, _ := util.ParseToken(c.GetHeader("Authorization"))
	if err := c.ShouldBind(&service); err == nil {
		res := service.ChangePwd(c, claims.ID, oldpassword, newpassword)
		c.JSON(200, res)
	} else {
		c.JSON(400, err)
	}
}

// 创建用户
func CreateUser(c *gin.Context) {
	var service service.UserService
	claims, _ := util.ParseToken(c.GetHeader("Authorization"))
	if err := c.ShouldBind(&service); err == nil {
		res := service.CreateUser(c, claims.ID)
		c.JSON(200, res)
	} else {
		c.JSON(400, err)
	}
}

// 更新用户
func UpdateUserById(c *gin.Context) {
	var service service.UserService
	claims, _ := util.ParseToken(c.GetHeader("Authorization"))
	if err := c.ShouldBind(&service); err == nil {
		res := service.UpdateUserById(c, claims.ID)
		c.JSON(200, res)
	} else {
		c.JSON(400, err)
	}
}

// 批量删除用户
func BatchDeleteUserByIds(c *gin.Context) {
	var service service.UserService
	claims, _ := util.ParseToken(c.GetHeader("Authorization"))
	if err := c.ShouldBind(&service); err == nil {
		res := service.BatchDeleteUserByIds(c, claims.ID)
		c.JSON(200, res)
	} else {
		c.JSON(400, err)
	}
}
