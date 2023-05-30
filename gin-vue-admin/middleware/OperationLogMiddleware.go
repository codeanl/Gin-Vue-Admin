package middleware

import (
	"gin-vue-admin/gin-vue-admin/conf"
	"gin-vue-admin/gin-vue-admin/dao"
	"gin-vue-admin/gin-vue-admin/model"
	"gin-vue-admin/gin-vue-admin/util"
	"github.com/gin-gonic/gin"
	"strings"
	"time"
)

// 操作日志channel
var OperationLogChan = make(chan *model.OperationLog, 30)

func OperationLogMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 开始时间
		startTime := time.Now()
		// 处理请求
		c.Next()
		// 结束时间
		endTime := time.Now()
		// 执行耗时
		timeCost := endTime.Sub(startTime).Milliseconds()
		// 获取当前登录用户
		var username string
		token := c.GetHeader("Authorization")
		if token == "" {
			username = "未登录"
		}
		claims, _ := util.ParseToken(token)
		user, _ := dao.GetUserById(claims.ID)
		username = user.Username
		// 获取访问路径
		path := strings.TrimPrefix(c.FullPath(), "/"+conf.Conf.System.UrlPathPrefix)
		// 请求方式
		method := c.Request.Method
		// 获取接口描述
		apiDesc, _ := dao.GetApiDescByPath(path, method)
		operationLog := model.OperationLog{
			Username:   username,
			Ip:         c.ClientIP(),
			IpLocation: "",
			Method:     method,
			Path:       path,
			Desc:       apiDesc,
			Status:     c.Writer.Status(),
			StartTime:  startTime,
			TimeCost:   timeCost,
		}
		// 最好是将日志发送到rabbitmq或者kafka中
		// 这里是发送到channel中，开启3个goroutine处理
		OperationLogChan <- &operationLog
	}
}
