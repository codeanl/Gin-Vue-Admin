package dao

import (
	"gin-vue-admin/gin-vue-admin/common"
	"gin-vue-admin/gin-vue-admin/model"
)

// 获取菜单列表
func GetMenus() ([]*model.Menu, error) {
	var menus []*model.Menu
	err := common.DB.Order("sort").Find(&menus).Error
	return menus, err
}

// 获取菜单树
func GetMenuTree() ([]*model.Menu, error) {
	var menus []*model.Menu
	err := common.DB.Order("sort").Find(&menus).Error
	// parentId为0的是根菜单
	return GenMenuTree(0, menus), err
}
func GenMenuTree(parentId uint, menus []*model.Menu) []*model.Menu {
	tree := make([]*model.Menu, 0)

	for _, m := range menus {
		if *m.ParentId == parentId {
			children := GenMenuTree(m.ID, menus)
			m.Children = children
			tree = append(tree, m)
		}
	}
	return tree
}

// 创建菜单
func CreateMenu(menu *model.Menu) error {
	err := common.DB.Create(menu).Error
	return err
}

// 更新菜单
func UpdateMenuById(menuId uint, menu *model.Menu) error {
	err := common.DB.Model(menu).Where("id = ?", menuId).Updates(menu).Error
	return err
}

// 批量删除菜单
func BatchDeleteMenuByIds(menuIds []uint) error {
	var menus []*model.Menu
	err := common.DB.Where("id IN (?)", menuIds).Find(&menus).Error
	if err != nil {
		return err
	}
	err = common.DB.Select("Roles").Unscoped().Delete(&menus).Error
	return err
}

// 根据用户ID获取用户的权限(可访问)菜单树
func GetUserMenuTreeByUserId(userId uint) ([]*model.Menu, error) {
	menus, err := GetUserMenusByUserId(userId)
	if err != nil {
		return nil, err
	}
	tree := GenMenuTree(0, menus)
	return tree, err
}
