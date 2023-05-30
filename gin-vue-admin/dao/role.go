package dao

import (
	"errors"
	"fmt"
	"gin-vue-admin/gin-vue-admin/common"
	"gin-vue-admin/gin-vue-admin/model"
	"github.com/gin-gonic/gin"
	"github.com/thoas/go-funk"
)

// 获取角色列表
func GetRoles(name, keyword string, status, pageNum, pageSize int) ([]model.Role, int64, error) {
	var list []model.Role
	db := common.DB.Model(&model.Role{}).Order("created_at DESC")

	if name != "" {
		db = db.Where("name LIKE ?", fmt.Sprintf("%%%s%%", name))
	}
	if keyword != "" {
		db = db.Where("keyword LIKE ?", fmt.Sprintf("%%%s%%", keyword))
	}
	if status != 0 {
		db = db.Where("status = ?", status)
	}
	// 当pageNum > 0 且 pageSize > 0 才分页
	//记录总条数
	var total int64
	err := db.Count(&total).Error
	if err != nil {
		return list, total, err
	}
	if pageNum > 0 && pageSize > 0 {
		err = db.Offset((pageNum - 1) * pageSize).Limit(pageSize).Find(&list).Error
	} else {
		err = db.Find(&list).Error
	}
	return list, total, err
}

// 查询用户是否存在
func CheckRoleExist(name string) (role *model.Role, exist bool, err error) {
	var count int64
	err = common.DB.Model(&model.Role{}).Where("name=?", name).
		Count(&count).Error
	if count == 0 {
		return nil, false, err
	}
	err = common.DB.Model(&model.Role{}).Where("name=?", name).
		First(&role).Error
	if err != nil {
		return nil, false, err
	}
	return role, true, nil
}

// 获取当前用户角色排序最小值（最高等级角色）以及当前用户信息
func GetCurrentUserMinRoleSort(c *gin.Context, uid uint) (uint, model.User, error) {
	// 获取当前用户
	ctxUser, err := GetUserById(uid)
	if err != nil {
		return 999, ctxUser, err
	}
	// 获取当前用户的所有角色
	currentRoles := ctxUser.Roles
	// 获取当前用户角色的排序，和前端传来的角色排序做比较
	var currentRoleSorts []int
	for _, role := range currentRoles {
		currentRoleSorts = append(currentRoleSorts, int(role.Sort))
	}
	// 当前用户角色排序最小值（最高等级角色）
	currentRoleSortMin := uint(funk.MinInt(currentRoleSorts))
	return currentRoleSortMin, ctxUser, nil
}

// 创建角色
func CreateRole(role *model.Role) error {
	err := common.DB.Create(role).Error
	return err
}

// 根据角色ID获取角色
func GetRolesByIds(roleIds []uint) ([]*model.Role, error) {
	var list []*model.Role
	err := common.DB.Where("id IN (?)", roleIds).Find(&list).Error
	return list, err
}

// 更新角色
func UpdateRoleById(roleId uint, role *model.Role) error {
	err := common.DB.Model(&model.Role{}).Where("id = ?", roleId).Updates(role).Error
	return err
}

// 获取角色的权限菜单
func GetRoleMenusById(roleId uint) ([]*model.Menu, error) {
	var role model.Role
	err := common.DB.Where("id = ?", roleId).Preload("Menus").First(&role).Error
	return role.Menus, err
}

// 根据用户ID获取用户的权限(可访问)菜单列表
func GetUserMenusByUserId(userId uint) ([]*model.Menu, error) {
	// 获取用户
	var user model.User
	err := common.DB.Where("id = ?", userId).Preload("Roles").First(&user).Error
	if err != nil {
		return nil, err
	}
	// 获取角色
	roles := user.Roles
	// 所有角色的菜单集合
	allRoleMenus := make([]*model.Menu, 0)
	for _, role := range roles {
		var userRole model.Role
		err := common.DB.Where("id = ?", role.ID).Preload("Menus").First(&userRole).Error
		if err != nil {
			return nil, err
		}
		// 获取角色的菜单
		menus := userRole.Menus
		allRoleMenus = append(allRoleMenus, menus...)
	}

	// 所有角色的菜单集合去重
	allRoleMenusId := make([]int, 0)
	for _, menu := range allRoleMenus {
		allRoleMenusId = append(allRoleMenusId, int(menu.ID))
	}
	allRoleMenusIdUniq := funk.UniqInt(allRoleMenusId)
	allRoleMenusUniq := make([]*model.Menu, 0)
	for _, id := range allRoleMenusIdUniq {
		for _, menu := range allRoleMenus {
			if id == int(menu.ID) {
				allRoleMenusUniq = append(allRoleMenusUniq, menu)
				break
			}
		}
	}

	// 获取状态status为1的菜单
	accessMenus := make([]*model.Menu, 0)
	for _, menu := range allRoleMenusUniq {
		if menu.Status == 1 {
			accessMenus = append(accessMenus, menu)
		}
	}

	return accessMenus, err
}

// 更新角色的权限菜单
func UpdateRoleMenus(role *model.Role) error {
	err := common.DB.Model(role).Association("Menus").Replace(role.Menus)
	return err
}

// 更新角色的权限接口（先全部删除再新增）
func UpdateRoleApis(roleKeyword string, reqRolePolicies [][]string) error {
	// 先获取path中的角色ID对应角色已有的police(需要先删除的)
	err := common.CasbinEnforcer.LoadPolicy()
	if err != nil {
		return errors.New("角色的权限接口策略加载失败")
	}
	rmPolicies := common.CasbinEnforcer.GetFilteredPolicy(0, roleKeyword)
	if len(rmPolicies) > 0 {
		isRemoved, _ := common.CasbinEnforcer.RemovePolicies(rmPolicies)
		if !isRemoved {
			return errors.New("更新角色的权限接口失败")
		}
	}
	isAdded, _ := common.CasbinEnforcer.AddPolicies(reqRolePolicies)
	if !isAdded {
		return errors.New("更新角色的权限接口失败")
	}
	err = common.CasbinEnforcer.LoadPolicy()
	if err != nil {
		return errors.New("更新角色的权限接口成功，角色的权限接口策略加载失败")
	} else {
		return err
	}
}

// 删除角色
func BatchDeleteRoleByIds(roleIds []uint) error {
	var roles []*model.Role
	err := common.DB.Where("id IN (?)", roleIds).Find(&roles).Error
	if err != nil {
		return err
	}
	err = common.DB.Select("Users", "Menus").Unscoped().Delete(&roles).Error
	// 删除成功就删除casbin policy
	if err == nil {
		for _, role := range roles {
			roleKeyword := role.Keyword
			rmPolicies := common.CasbinEnforcer.GetFilteredPolicy(0, roleKeyword)
			if len(rmPolicies) > 0 {
				isRemoved, _ := common.CasbinEnforcer.RemovePolicies(rmPolicies)
				if !isRemoved {
					return errors.New("删除角色成功, 删除角色关联权限接口失败")
				}
			}
		}

	}
	return err
}

// 根据角色关键字获取角色的权限接口
func GetRoleApisByRoleKeyword(roleKeyword string) ([]*model.Api, error) {
	policies := common.CasbinEnforcer.GetFilteredPolicy(0, roleKeyword)

	// 获取所有接口
	var apis []*model.Api
	err := common.DB.Find(&apis).Error
	if err != nil {
		return apis, errors.New("获取角色的权限接口失败")
	}

	accessApis := make([]*model.Api, 0)

	for _, policy := range policies {
		path := policy[1]
		method := policy[2]
		for _, api := range apis {
			if path == api.Path && method == api.Method {
				accessApis = append(accessApis, api)
				break
			}
		}
	}

	return accessApis, err
}
