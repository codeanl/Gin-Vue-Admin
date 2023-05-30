package dao

import (
	"errors"
	"fmt"
	"gin-vue-admin/gin-vue-admin/common"
	"gin-vue-admin/gin-vue-admin/model"
	"github.com/thoas/go-funk"
)

// 获取单个用户
func GetUserById(id uint) (model.User, error) {
	var user model.User
	err := common.DB.Where("id = ?", id).Preload("Roles").First(&user).Error
	return user, err
}

// 根据username判断是否存在该名字
func ExistOrNotByUserName(userName string) (user *model.User, exist bool, err error) {
	var count int64
	err = common.DB.Model(&model.User{}).Where("username=?", userName).
		Count(&count).Error
	if count == 0 {
		return nil, false, err
	}
	err = common.DB.Model(&model.User{}).Where("username=?", userName).
		First(&user).Error
	if err != nil {
		return nil, false, err
	}
	return user, true, nil
}

// 获取用户列表
func GetUsers(username, nickname, mobile string, status uint, pageNum, pageSize int) ([]*model.User, int64, error) {
	var list []*model.User
	db := common.DB.Model(&model.User{}).Order("created_at DESC")

	if username != "" {
		db = db.Where("username LIKE ?", fmt.Sprintf("%%%s%%", username))
	}
	if nickname != "" {
		db = db.Where("nickname LIKE ?", fmt.Sprintf("%%%s%%", nickname))
	}
	if mobile != "" {
		db = db.Where("mobile LIKE ?", fmt.Sprintf("%%%s%%", mobile))
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
		err = db.Offset((pageNum - 1) * pageSize).Limit(pageSize).Preload("Roles").Find(&list).Error
	} else {
		err = db.Preload("Roles").Find(&list).Error
	}
	return list, total, err
}

// 更新用户信息
func UpdateUserById(uId uint, user model.User) error {
	return common.DB.Model(&model.User{}).Where("id=?", uId).
		Updates(&user).Error
}

// 创建用户
func CreateUser(user *model.User) error {
	err := common.DB.Create(user).Error
	return err
}

// 更新用户
func UpdateUser(user *model.User) error {
	err := common.DB.Model(user).Updates(user).Error
	if err != nil {
		return err
	}
	err = common.DB.Model(user).Association("Roles").Replace(user.Roles)

	return err
}

// 根据用户ID获取用户角色排序最小值
func GetUserMinRoleSortsByIds(ids []uint) ([]int, error) {
	// 根据用户ID获取用户信息
	var userList []model.User
	err := common.DB.Where("id IN (?)", ids).Preload("Roles").Find(&userList).Error
	if err != nil {
		return []int{}, err
	}
	if len(userList) == 0 {
		return []int{}, errors.New("未获取到任何用户信息")
	}
	var roleMinSortList []int
	for _, user := range userList {
		roles := user.Roles
		var roleSortList []int
		for _, role := range roles {
			roleSortList = append(roleSortList, int(role.Sort))
		}
		roleMinSort := funk.MinInt(roleSortList)
		roleMinSortList = append(roleMinSortList, roleMinSort)
	}
	return roleMinSortList, nil
}

// 批量删除
func BatchDeleteUserByIds(ids []uint) error {
	// 用户和角色存在多对多关联关系
	var users []model.User
	for _, id := range ids {
		// 根据ID获取用户
		user, err := GetUserById(id)
		if err != nil {
			return errors.New(fmt.Sprintf("未获取到ID为%d的用户", id))
		}
		users = append(users, user)
	}

	err := common.DB.Select("Roles").Unscoped().Delete(&users).Error
	return err
}
