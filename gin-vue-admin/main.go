package main

import (
	"gin-vue-admin/gin-vue-admin/common"
	"gin-vue-admin/gin-vue-admin/conf"
	"gin-vue-admin/gin-vue-admin/router"
)

func main() {
	// 加载配置文件到全局配置结构体
	conf.InitConfig()
	// 初始化日志
	common.InitLogger()
	// 初始化数据库(mysql)
	common.InitMysql()
	// 初始化casbin策略管理器
	common.InitCasbinEnforcer()
	// 注册所有路由
	router.InitRoutes()
}
