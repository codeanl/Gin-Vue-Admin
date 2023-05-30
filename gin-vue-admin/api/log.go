package api

import (
	"gin-vue-admin/gin-vue-admin/service"
	"github.com/gin-gonic/gin"
)

// 获取操作日志列表
func GetOperationLogs(c *gin.Context) {
	var service service.OperationLogService
	if err := c.ShouldBind(&service); err == nil {
		res := service.GetOperationLogs(c)
		c.JSON(200, res)
	} else {
		c.JSON(400, err)
	}
}

// 批量删除操作日志
func BatchDeleteOperationLogByIds(c *gin.Context) {
	var service service.OperationLogService
	if err := c.ShouldBind(&service); err == nil {
		res := service.BatchDeleteOperationLogByIds(c)
		c.JSON(200, res)
	} else {
		c.JSON(400, err)
	}
}
