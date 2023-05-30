package router

import (
	"gin-vue-admin/gin-vue-admin/conf"
	"gin-vue-admin/gin-vue-admin/middleware"
	"github.com/gin-gonic/gin"
	"log"
	"time"
)

// 初始化
func InitRoutes() *gin.Engine {
	//设置模式
	gin.SetMode(conf.Conf.System.Mode)

	// 创建带有默认中间件的路由:
	// 日志与恢复中间件
	r := gin.Default()

	// 启用限流中间件
	// 默认每50毫秒填充一个令牌，最多填充200个
	fillInterval := time.Duration(conf.Conf.RateLimit.FillInterval)
	capacity := conf.Conf.RateLimit.Capacity
	r.Use(middleware.RateLimitMiddleware(time.Millisecond*fillInterval, capacity))

	// 启用全局跨域中间件
	r.Use(middleware.CORSMiddleware())

	// 启用操作日志中间件
	r.Use(middleware.OperationLogMiddleware())

	// 路由分组
	apiGroup := r.Group("/" + conf.Conf.System.UrlPathPrefix)
	//
	//注册路由
	InitBaseRoutes(apiGroup)         // 注册基础路由, 不需要jwt认证中间件,不需要casbin中间件
	InitUserRoutes(apiGroup)         // 注册用户路由, jwt认证中间件,casbin鉴权中间件
	InitRoleRoutes(apiGroup)         // 注册角色路由, jwt认证中间件,casbin鉴权中间件
	InitMenuRoutes(apiGroup)         // 注册菜单路由, jwt认证中间件,casbin鉴权中间件
	InitApiRoutes(apiGroup)          // 注册接口路由, jwt认证中间件,casbin鉴权中间件
	InitOperationLogRoutes(apiGroup) // 注册操作日志路由, jwt认证中间件,casbin鉴权中间件
	log.Println(`
		 ┏┓　　　┏┓
		┏┛┻━━━━━┛┻┓
		┃　　      ┃
		┃　　　━　　┃
		┃　┳┛　┗┳  ┃
		┃　　　　　 ┃
		┃　　　┻　　┃
		┃　　　　　 ┃
		┗━┓　　　┏━┛
		  ┃　　　┃   神兽保佑　　　　　　　　
		  ┃　　　┃   代码无BUG！
		  ┃　　　┗━━━┓
		  ┃　　　　　 ┣┓
		  ┃　　　　　 ┏┛
		  ┗┓┓┏━┳┓  ┏┛
		    ┃┫┫　┃┫┫
		    ┗┻┛　┗┻┛
		`)
	r.Run(conf.Conf.System.Port)
	return r
}
