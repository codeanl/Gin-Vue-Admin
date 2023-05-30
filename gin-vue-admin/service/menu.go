package service

import (
	"gin-vue-admin/gin-vue-admin/dao"
	"gin-vue-admin/gin-vue-admin/model"
	"gin-vue-admin/gin-vue-admin/response"
	"gin-vue-admin/gin-vue-admin/util/e"
	"github.com/gin-gonic/gin"
)

type MenuService struct {
	ID         uint          `gorm:"type:varchar(20)" json:"id" form:"id"`
	UserID     uint          `gorm:"type:varchar(20)" json:"userId" form:"userId"`
	Name       string        `gorm:"type:varchar(50);comment:'菜单名称(英文名, 可用于国际化)'" json:"name" form:"name"`
	Title      string        `gorm:"type:varchar(50);comment:'菜单标题(无法国际化时使用)'" json:"title" form:"title"`
	Icon       string        `gorm:"type:varchar(50);comment:'菜单图标'" json:"icon" form:"icon"`
	Path       string        `gorm:"type:varchar(100);comment:'菜单访问路径'" json:"path" form:"path"`
	Redirect   string        `gorm:"type:varchar(100);comment:'重定向路径'" json:"redirect" form:"redirect"`
	Component  string        `gorm:"type:varchar(100);comment:'前端组件路径'" json:"component" form:"component"`
	Sort       uint          `gorm:"type:int(3) unsigned;default:999;comment:'菜单顺序(1-999)'" json:"sort" form:"sort"`
	Status     uint          `gorm:"type:tinyint(1);default:1;comment:'菜单状态(正常/禁用, 默认正常)'" json:"status" form:"status"`
	Hidden     uint          `gorm:"type:tinyint(1);default:2;comment:'菜单在侧边栏隐藏(1隐藏，2显示)'" json:"hidden" form:"hidden"`
	NoCache    uint          `gorm:"type:tinyint(1);default:2;comment:'菜单是否被 <keep-alive> 缓存(1不缓存，2缓存)'" json:"noCache" form:"noCache"`
	AlwaysShow uint          `gorm:"type:tinyint(1);default:2;comment:'忽略之前定义的规则，一直显示根路由(1忽略，2不忽略)'" json:"alwaysShow" form:"alwaysShow"`
	Breadcrumb uint          `gorm:"type:tinyint(1);default:1;comment:'面包屑可见性(可见/隐藏, 默认可见)'" json:"breadcrumb" form:"breadcrumb"`
	ActiveMenu *string       `gorm:"type:varchar(100);comment:'在其它路由时，想在侧边栏高亮的路由'" json:"activeMenu" form:"activeMenu"`
	ParentId   *uint         `gorm:"default:0;comment:'父菜单编号(编号为0时表示根菜单)'" json:"parentId" form:"parentId"`
	Creator    string        `gorm:"type:varchar(20);comment:'创建人'" json:"creator" form:"creator"`
	Children   []*model.Menu `gorm:"-" json:"children"`                  // 子菜单集合
	Roles      []*model.Role `gorm:"many2many:role_menus;" json:"roles"` // 角色菜单多对多关系
	MenuIds    []uint        `json:"menuIds" form:"menuIds"`
}

// GetMenus 获取菜单列表
func (req MenuService) GetMenus(c *gin.Context) response.Res {
	code := e.SUCCESS
	menus, err := dao.GetMenus()
	if err != nil {
		code = e.ERROR
		return response.Res{
			Code: code,
			Msg:  e.GetMsg(code),
		}
	}
	return response.Res{
		Code: code,
		Msg:  e.GetMsg(code),
		Data: menus,
	}
}

// GetMenuTree 获取菜单树
func (req MenuService) GetMenuTree(c *gin.Context) response.Res {
	code := e.SUCCESS

	menuTree, err := dao.GetMenuTree()
	if err != nil {
		code = e.ERROR
		return response.Res{
			Code: code,
			Msg:  e.GetMsg(code),
		}
	}
	return response.Res{
		Code: code,
		Msg:  e.GetMsg(code),
		Data: menuTree,
	}
}

// CreateMenu 创建菜单
func (req MenuService) CreateMenu(c *gin.Context, myid uint) response.Res {
	code := e.SUCCESS
	// 获取当前用户
	ctxUser, err := dao.GetUserById(myid)
	if err != nil {
		code = e.ERROR
		return response.Res{
			Code: code,
			Msg:  e.GetMsg(code),
		}
	}
	menu := model.Menu{
		Name:       req.Name,
		Title:      req.Title,
		Icon:       req.Icon,
		Path:       req.Path,
		Redirect:   req.Redirect,
		Component:  req.Component,
		Sort:       req.Sort,
		Status:     req.Status,
		Hidden:     req.Hidden,
		NoCache:    req.NoCache,
		AlwaysShow: req.AlwaysShow,
		Breadcrumb: req.Breadcrumb,
		ActiveMenu: req.ActiveMenu,
		ParentId:   req.ParentId,
		Creator:    ctxUser.Username,
	}
	err = dao.CreateMenu(&menu)
	if err != nil {
		code = e.ERROR
		return response.Res{
			Code: code,
			Msg:  e.GetMsg(code),
		}
	}
	return response.Res{
		Code: code,
		Msg:  e.GetMsg(code),
		Data: menu,
	}
}

// UpdateMenuById 更新菜单
func (req MenuService) UpdateMenuById(c *gin.Context, myid uint) response.Res {
	code := e.SUCCESS
	// 获取路径中的menuId
	if req.ID <= 0 {
		code = e.ApiIDError
		return response.Res{
			Code: code,
			Msg:  e.GetMsg(code),
		}
	}
	// 获取当前用户
	ctxUser, err := dao.GetUserById(myid)
	if err != nil {
		code = e.ERROR
		return response.Res{
			Code: code,
			Msg:  e.GetMsg(code),
		}
	}
	menu := model.Menu{
		Name:       req.Name,
		Title:      req.Title,
		Icon:       req.Icon,
		Path:       req.Path,
		Redirect:   req.Redirect,
		Component:  req.Component,
		Sort:       req.Sort,
		Status:     req.Status,
		Hidden:     req.Hidden,
		NoCache:    req.NoCache,
		AlwaysShow: req.AlwaysShow,
		Breadcrumb: req.Breadcrumb,
		ActiveMenu: req.ActiveMenu,
		ParentId:   req.ParentId,
		Creator:    ctxUser.Username,
	}
	err = dao.UpdateMenuById(req.ID, &menu)
	if err != nil {
		code = e.ERROR
		return response.Res{
			Code: code,
			Msg:  e.GetMsg(code),
		}
	}
	return response.Res{
		Code: code,
		Msg:  e.GetMsg(code),
		Data: menu,
	}

}

// BatchDeleteMenuByIds 批量删除菜单
func (req MenuService) BatchDeleteMenuByIds(c *gin.Context) response.Res {
	code := e.SUCCESS
	err := dao.BatchDeleteMenuByIds(req.MenuIds)
	if err != nil {
		code = e.ERROR
		return response.Res{
			Code: code,
			Msg:  e.GetMsg(code),
		}
	}
	return response.Res{
		Code: code,
		Msg:  e.GetMsg(code),
	}
}

// GetUserMenusByUserId 根据用户ID获取用户的可访问菜单列表
func (req MenuService) GetUserMenusByUserId(c *gin.Context) response.Res {
	code := e.SUCCESS
	// 获取路径中的userId
	if req.UserID <= 0 {
		code = e.UserIDError
		return response.Res{
			Code: code,
			Msg:  e.GetMsg(code),
		}
	}
	menus, err := dao.GetUserMenusByUserId(req.UserID)
	if err != nil {
		code = e.ERROR
		return response.Res{
			Code: code,
			Msg:  e.GetMsg(code),
		}
	}
	return response.Res{
		Code: code,
		Msg:  e.GetMsg(code),
		Data: menus,
	}
}

// 根据用户ID获取用户的可访问菜单树
func (req MenuService) GetUserMenuTreeByUserId(c *gin.Context) response.Res {
	code := e.SUCCESS
	// 获取路径中的userId
	if req.UserID <= 0 {
		code = e.UserIDError
		return response.Res{
			Code: code,
			Msg:  e.GetMsg(code),
		}
	}
	menus, err := dao.GetUserMenuTreeByUserId(req.UserID)
	if err != nil {
		code = e.ERROR
		return response.Res{
			Code: code,
			Msg:  e.GetMsg(code),
		}
	}
	return response.Res{
		Code: code,
		Msg:  e.GetMsg(code),
		Data: menus,
	}
}
