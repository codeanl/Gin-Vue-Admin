<h1 align="center">Gin-Vue-Admin后台管理脚手架</h1>

<div align="center">
Gin+Vue开发的管理系统脚手架,前后端分离,基于角色的访问控制(RBAC),易于扩展。 后端Go包含了gin、 gorm、 jwt和casbin等的使用, 前端用vue（未完成）。
<p align="center">
<img src="https://img.shields.io/badge/Gin-1.9.0-brightgreen" alt="Gin version"/>
<img src="https://img.shields.io/badge/Gorm-1.25.1-brightgreen" alt="Gorm version"/>
<img src="https://img.shields.io/github/license/gnimli/go-web-mini" alt="License"/>
</p>
</div>

## 简介

Gin-Vue-Admin后台管理脚手架是一款轻巧的后台管理脚手架, 专为普通用户设计，开箱即用，无需复杂配置。我们的目标是打造最轻量化的后台管理系统！


## 项目结构概览

```
├─common      #公共资源
├─conf        # 配置文件
├─api         # api层响应路由请求
├─dao         # 数据库操作
├─middleware  # 中间件
├─model       # 结构体模型
├─logs        # 日志返回
├─response    # 常用返回封装
├─router      # 所有路由
├─util        # 工具方法
└─service     # 逻辑处理

```

## 模块进度
### 用户模块
- [x] 获取用户列表
- [x] 更新用户登录密码
- [x] 创建用户
- [x] 更新用户
- [x] 批量删除用户
- [x] 用户登录
- [x] 退出登录
- [x] 获取自己的详细信息

### 角色模块
- [x] 获取角色列表
- [x] 创建角色
- [x] 更新角色
- [x] 获取角色的权限菜单
- [x] 更新角色的权限菜单
- [x] 根据角色关键字获取角色的权限接口
- [x] 更新角色的权限接口
- [x] 批量删除角色

### 菜单模块
- [x] 获取菜单列表
- [x] 获取菜单树
- [x] 创建菜单
- [x] 更新菜单
- [x] 根据用户ID获取用户的可访问菜单列表
- [x] 根据用户ID获取用户的可访问菜单树

### 接口模块
- [x] 获取接口列表
- [x] 获取接口树
- [x] 创建接口
- [x] 更新接口
- [x] 批量删除接口

###日志模块
- [x] 获取操作日志列表
- [x] 批量删除操作日志

## 安装
```
后端
cd gin-vue-admin
go mod tidy
go run main.go
```


## TODO
- [ ] 完成没有操作日志的bug
- [ ] 完成前端

