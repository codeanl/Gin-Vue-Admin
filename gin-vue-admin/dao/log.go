package dao

import (
	"fmt"
	"gin-vue-admin/gin-vue-admin/common"
	"gin-vue-admin/gin-vue-admin/model"
)

func GetOperationLogs(username, ip, path string, status, pageNum, pageSize int) ([]model.OperationLog, int64, error) {
	var list []model.OperationLog
	db := common.DB.Model(&model.OperationLog{}).Order("start_time DESC")
	if username != "" {
		db = db.Where("username LIKE ?", fmt.Sprintf("%%%s%%", username))
	}
	if ip != "" {
		db = db.Where("ip LIKE ?", fmt.Sprintf("%%%s%%", ip))
	}
	if path != "" {
		db = db.Where("path LIKE ?", fmt.Sprintf("%%%s%%", path))
	}
	if status != 0 {
		db = db.Where("status = ?", status)
	}
	// 分页
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
func BatchDeleteOperationLogByIds(ids []uint) error {
	err := common.DB.Where("id IN (?)", ids).Unscoped().Delete(&model.OperationLog{}).Error
	return err
}

// 根据接口路径和请求方式获取接口描述
func GetApiDescByPath(path string, method string) (string, error) {
	var api model.Api
	err := common.DB.Where("path = ?", path).Where("method = ?", method).First(&api).Error
	return api.Desc, err
}
