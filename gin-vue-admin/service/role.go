package service

import (
	"fmt"
	"gin-vue-admin/gin-vue-admin/common"
	"gin-vue-admin/gin-vue-admin/dao"
	"gin-vue-admin/gin-vue-admin/model"
	"gin-vue-admin/gin-vue-admin/response"
	"gin-vue-admin/gin-vue-admin/util/e"
	"github.com/gin-gonic/gin"
	"github.com/thoas/go-funk"
)

type RoleService struct {
	ID       uint          `gorm:"type:varchar(20)" json:"id" form:"id"`
	Name     string        `gorm:"type:varchar(20);not null;unique" json:"name" form:"name"`
	Keyword  string        `gorm:"type:varchar(20);not null;unique" json:"keyword" form:"keyword"`
	Desc     string        `gorm:"type:varchar(100);" json:"desc" form:"desc"`
	Status   uint          `gorm:"type:tinyint(1);default:1;comment:'1正常, 2禁用'" json:"status" form:"status"`
	Sort     uint          `gorm:"type:int(3);default:999;comment:'角色排序(排序越大权限越低, 不能查看比自己序号小的角色, 不能编辑同序号用户权限, 排序为1表示超级管理员)'" json:"sort" form:"sort"`
	Creator  string        `gorm:"type:varchar(20);" json:"creator" form:"creator"`
	Users    []*model.User `gorm:"many2many:user_roles" json:"users" form:"users"`
	Menus    []*model.Menu `gorm:"many2many:role_menus;" json:"menus"` // 角色菜单多对多关系
	PageNum  uint          `json:"pageNum" form:"pageNum"`
	PageSize uint          `json:"pageSize" form:"pageSize"`
	MenuIds  []uint        `json:"menuIds" form:"menuIds"`
	ApiIds   []uint        `json:"apiIds" form:"apiIds"`
	RoleIds  []uint        `json:"roleIds" form:"roleIds"`
}

// GetRoles 获取角色列表
func (req RoleService) GetRoles(c *gin.Context) response.Res {
	code := e.SUCCESS
	// 获取角色列表
	roles, total, err := dao.GetRoles(req.Name, req.Keyword, int(req.Status), int(req.PageNum), int(req.PageSize))
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
		Data: map[string]interface{}{
			"total": total,
			"roles": roles,
		},
	}
}

// CreateRole 创建角色
func (req RoleService) CreateRole(c *gin.Context, uid uint) response.Res {
	code := e.SUCCESS
	if req.Name == "" || req.Keyword == "" || req.Desc == "" || req.Sort == 0 {
		code = e.InvalidParams
		return response.Res{
			Code: code,
			Msg:  e.GetMsg(code),
		}
	}
	_, exist, err := dao.CheckRoleExist(req.Name)
	if exist {
		code = e.RoleNameExist
		return response.Res{
			Code: code,
			Msg:  e.GetMsg(code),
		}
	}
	// 获取当前用户最高角色等级
	sort, ctxUser, err := dao.GetCurrentUserMinRoleSort(c, uid)
	if err != nil {
		code = e.ERROR
		return response.Res{
			Code: code,
			Msg:  e.GetMsg(code),
		}
	}
	// 用户不能创建比自己等级高或相同等级的角色
	if sort >= req.Sort {
		code = e.CreatRoleGreater
		return response.Res{
			Code: code,
			Msg:  e.GetMsg(code),
		}
	}
	role := model.Role{
		Name:    req.Name,
		Keyword: req.Keyword,
		Desc:    req.Desc,
		Status:  req.Status,
		Sort:    req.Sort,
		Creator: ctxUser.Username,
	}
	// 创建角色
	err = dao.CreateRole(&role)
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
		Data: role,
	}
}

// UpdateRoleById 更新角色
func (req RoleService) UpdateRoleById(c *gin.Context, myid uint) response.Res {
	code := e.SUCCESS
	_, exist, err := dao.CheckRoleExist(req.Name)
	if exist {
		code = e.RoleNameExist
		return response.Res{
			Code: code,
			Msg:  e.GetMsg(code),
		}
	}
	if req.ID <= 0 {
		code = e.RoleIDError
		return response.Res{
			Code: code,
			Msg:  e.GetMsg(code),
		}
	}
	// 当前用户角色排序最小值（最高等级角色）以及当前用户
	minSort, ctxUser, err := dao.GetCurrentUserMinRoleSort(c, myid)
	if err != nil {
		code = e.ERROR
		return response.Res{
			Code: code,
			Msg:  e.GetMsg(code),
			Data: "!",
		}
	}
	// 不能更新比自己角色等级高或相等的角色
	// 根据path中的角色ID获取该角色信息
	roles, err := dao.GetRolesByIds([]uint{req.ID})
	if err != nil {
		code = e.ERROR
		return response.Res{
			Code: code,
			Msg:  e.GetMsg(code),
		}
	}
	if len(roles) == 0 {
		code = e.UserRoleError
		return response.Res{
			Code: code,
			Msg:  e.GetMsg(code),
		}
	}
	if minSort >= roles[0].Sort {
		code = e.RoleGreaterMe
		return response.Res{
			Code: code,
			Msg:  e.GetMsg(code),
		}
	}

	// 不能把角色等级更新得比当前用户的等级高
	if minSort >= req.Sort {
		code = e.RoleGreaterMe
		return response.Res{
			Code: code,
			Msg:  e.GetMsg(code),
		}
	}
	role := model.Role{
		Name:    req.Name,
		Keyword: req.Keyword,
		Desc:    req.Desc,
		Status:  req.Status,
		Sort:    req.Sort,
		Creator: ctxUser.Username,
	}
	// 更新角色
	err = dao.UpdateRoleById(req.ID, &role)
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
		Data: role,
	}
}

// GetRoleMenusById 获取角色的权限菜单
func (req RoleService) GetRoleMenusById(c *gin.Context) response.Res {
	code := e.SUCCESS
	// 获取path中的roleId
	if req.ID <= 0 {
		code = e.ERROR
		return response.Res{
			Code:  code,
			Msg:   e.GetMsg(code),
			Error: "角色ID不正确",
		}
	}
	menus, err := dao.GetRoleMenusById(req.ID)
	if err != nil {
		code = e.ERROR
		return response.Res{
			Code:  code,
			Msg:   e.GetMsg(code),
			Error: "获取角色的权限菜单失败" + err.Error(),
		}
	}
	return response.Res{
		Code: code,
		Msg:  e.GetMsg(code),
		Data: menus,
	}
}

// UpdateRoleMenusById 更新角色的权限菜单
func (req RoleService) UpdateRoleMenusById(c *gin.Context, myid uint) response.Res {
	code := e.SUCCESS

	// 获取path中的roleId
	if req.ID <= 0 {
		code = e.ERROR
		return response.Res{
			Code:  code,
			Msg:   e.GetMsg(code),
			Error: "角色ID不正确",
		}
	}
	// 根据path中的角色ID获取该角色信息
	roles, err := dao.GetRolesByIds([]uint{req.ID})
	if err != nil {
		code = e.ERROR
		return response.Res{
			Code: code,
			Msg:  e.GetMsg(code),
		}
	}
	if len(roles) == 0 {
		code = e.ERROR
		return response.Res{
			Code:  code,
			Msg:   e.GetMsg(code),
			Error: "未获取到角色信息",
		}
	}

	// 当前用户角色排序最小值（最高等级角色）以及当前用户
	minSort, ctxUser, err := dao.GetCurrentUserMinRoleSort(c, myid)
	if err != nil {
		code = e.ERROR
		return response.Res{
			Code: code,
			Msg:  e.GetMsg(code),
		}
	}

	// (非管理员)不能更新比自己角色等级高或相等角色的权限菜单
	if minSort != 1 {
		if minSort >= roles[0].Sort {
			code = e.ERROR
			return response.Res{
				Code:  code,
				Msg:   e.GetMsg(code),
				Error: "不能更新比自己角色等级高或相等角色的权限菜单",
			}
		}
	}

	// 获取当前用户所拥有的权限菜单
	ctxUserMenus, err := dao.GetUserMenusByUserId(ctxUser.ID)
	if err != nil {
		code = e.ERROR
		return response.Res{
			Code:  code,
			Msg:   e.GetMsg(code),
			Error: "获取当前用户的可访问菜单列表失败: " + err.Error(),
		}
	}

	// 获取当前用户所拥有的权限菜单ID
	ctxUserMenusIds := make([]uint, 0)
	for _, menu := range ctxUserMenus {
		ctxUserMenusIds = append(ctxUserMenusIds, menu.ID)
	}

	// 用户需要修改的菜单集合
	reqMenus := make([]*model.Menu, 0)
	menuIds := req.MenuIds
	// (非管理员)不能把角色的权限菜单设置的比当前用户所拥有的权限菜单多
	if minSort != 1 {
		for _, id := range menuIds {
			if !funk.Contains(ctxUserMenusIds, id) {
				code = e.ERROR
				return response.Res{
					Code:  code,
					Msg:   e.GetMsg(code),
					Error: fmt.Sprintf("无权设置ID为%d的菜单", id),
				}
			}
		}

		for _, id := range menuIds {
			for _, menu := range ctxUserMenus {
				if id == menu.ID {
					reqMenus = append(reqMenus, menu)
					break
				}
			}
		}
	} else {
		// 管理员随意设置
		// 根据menuIds查询查询菜单
		menus, err := dao.GetMenus()
		if err != nil {
			code := e.ERROR
			return response.Res{
				Code:  code,
				Msg:   e.GetMsg(code),
				Error: "获取菜单列表失败: ",
			}
		}
		for _, menuId := range menuIds {
			for _, menu := range menus {
				if menuId == menu.ID {
					reqMenus = append(reqMenus, menu)
				}
			}
		}
	}

	roles[0].Menus = reqMenus

	err = dao.UpdateRoleMenus(roles[0])
	if err != nil {
		code := e.ERROR
		return response.Res{
			Code:  code,
			Msg:   e.GetMsg(code),
			Error: "更新角色的权限菜单失败: " + err.Error(),
		}
	}
	return response.Res{
		Code: code,
		Msg:  e.GetMsg(code),
	}

}

// GetRoleApisById 获取角色的权限接口
func (req RoleService) GetRoleApisById(c *gin.Context) response.Res {
	code := e.SUCCESS
	// 获取path中的roleId
	if req.ID <= 0 {
		code = e.ERROR
		return response.Res{
			Code:  code,
			Msg:   e.GetMsg(code),
			Error: "角色ID不正确",
		}
	}
	// 根据path中的角色ID获取该角色信息
	roles, err := dao.GetRolesByIds([]uint{req.ID})
	if err != nil {
		code = e.ERROR
		return response.Res{
			Code: code,
			Msg:  e.GetMsg(code),
		}
	}
	if len(roles) == 0 {
		code = e.ERROR
		return response.Res{
			Code:  code,
			Msg:   e.GetMsg(code),
			Error: "未获取到角色信息",
		}
	}
	// 根据角色keyword获取casbin中policy
	keyword := roles[0].Keyword
	apis, err := dao.GetRoleApisByRoleKeyword(keyword)
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
		Data: apis,
	}
}

// UpdateRoleApisById 更新角色的权限接口
func (req RoleService) UpdateRoleApisById(c *gin.Context, myid uint) response.Res {
	code := e.SUCCESS
	// 获取path中的roleId
	if req.ID <= 0 {
		code = e.UserIDError
		return response.Res{
			Code: code,
			Msg:  e.GetMsg(code),
		}
	}
	// 根据path中的角色ID获取该角色信息
	roles, err := dao.GetRolesByIds([]uint{req.ID})
	if err != nil {
		code = e.ERROR
		return response.Res{
			Code: code,
			Msg:  e.GetMsg(code),
		}
	}
	if len(roles) == 0 {
		code = e.UserRoleError
		return response.Res{
			Code: code,
			Msg:  e.GetMsg(code),
		}
	}

	// 当前用户角色排序最小值（最高等级角色）以及当前用户
	minSort, ctxUser, err := dao.GetCurrentUserMinRoleSort(c, myid)
	if err != nil {
		code = e.ERROR
		return response.Res{
			Code: code,
			Msg:  e.GetMsg(code),
		}
	}
	// (非管理员)不能更新比自己角色等级高或相等角色的权限接口
	if minSort != 1 {
		if minSort >= roles[0].Sort {
			code = e.RoleGreaterMe
			return response.Res{
				Code: code,
				Msg:  e.GetMsg(code),
			}
		}
	}
	// 获取当前用户所拥有的权限接口
	ctxRoles := ctxUser.Roles
	ctxRolesPolicies := make([][]string, 0)
	for _, role := range ctxRoles {
		policy := common.CasbinEnforcer.GetFilteredPolicy(0, role.Keyword)
		ctxRolesPolicies = append(ctxRolesPolicies, policy...)
	}
	// 得到path中的角色ID对应角色能够设置的权限接口集合
	for _, policy := range ctxRolesPolicies {
		policy[0] = roles[0].Keyword
	}
	// 根据apiID获取接口详情
	apis, err := dao.GetApisById(req.ApiIds)
	if err != nil {
		code = e.ERROR
		return response.Res{
			Code: code,
			Msg:  e.GetMsg(code),
		}
	}
	// 生成前端想要设置的角色policies
	reqRolePolicies := make([][]string, 0)
	for _, api := range apis {
		reqRolePolicies = append(reqRolePolicies, []string{
			roles[0].Keyword, api.Path, api.Method,
		})
	}
	// (非管理员)不能把角色的权限接口设置的比当前用户所拥有的权限接口多
	if minSort != 1 {
		for _, reqPolicy := range reqRolePolicies {
			if !funk.Contains(ctxRolesPolicies, reqPolicy) {
				code = e.ERROR
				return response.Res{
					Code: code,
					Msg:  e.GetMsg(code),
				}
			}
		}
	}

	// 更新角色的权限接口
	err = dao.UpdateRoleApis(roles[0].Keyword, reqRolePolicies)
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

// BatchDeleteRoleByIds 批量删除角色
func (req RoleService) BatchDeleteRoleByIds(c *gin.Context, myid uint) response.Res {
	code := e.SUCCESS
	minSort, _, err := dao.GetCurrentUserMinRoleSort(c, myid)
	if err != nil {
		code = e.ERROR
		return response.Res{
			Code:  code,
			Msg:   e.GetMsg(code),
			Error: err.Error(),
		}
	}
	// 前端传来需要删除的角色ID
	// 获取角色信息
	roles, err := dao.GetRolesByIds(req.RoleIds)
	if err != nil {
		code = e.ERROR
		return response.Res{
			Code: code,
			Msg:  e.GetMsg(code),
		}
	}
	if len(roles) == 0 {
		code = e.UserRoleError
		return response.Res{
			Code: code,
			Msg:  e.GetMsg(code),
		}
	}
	// 不能删除比自己角色等级高或相等的角色
	for _, role := range roles {
		if minSort >= role.Sort {
			code = e.RoleGreaterMe
			return response.Res{
				Code: code,
				Msg:  e.GetMsg(code),
			}
		}
	}
	// 删除角色
	err = dao.BatchDeleteRoleByIds(req.RoleIds)
	if err != nil {
		code = e.ERROR
		return response.Res{
			Code: code,
			Msg:  e.GetMsg(code),
		}
	}
	// 删除角色成功直接清理缓存，让活跃的用户自己重新缓存最新用户信息
	return response.Res{
		Code: code,
		Msg:  e.GetMsg(code),
	}
}
