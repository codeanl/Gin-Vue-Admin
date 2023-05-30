package middleware

import (
	"gin-vue-admin/gin-vue-admin/util"
	"gin-vue-admin/gin-vue-admin/util/e"
	"github.com/gin-gonic/gin"
	"time"
)

// JWT token验证中间件
func InitAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		var code int
		var data interface{}
		code = 200
		token := c.GetHeader("Authorization")
		if token == "" {
			code = 404
		} else {
			claims, err := util.ParseToken(token)
			if err != nil {
				code = e.ERROR
			} else if time.Now().Unix() > claims.ExpiresAt {
				code = e.ERROR
			}
		}
		if code != e.SUCCESS {
			c.JSON(200, gin.H{
				"status": code,
				"msg":    e.GetMsg(code),
				"data":   data,
				"error":  "jwt验证失败",
			})
			c.Abort()
			return
		}
		c.Next()
	}
}
