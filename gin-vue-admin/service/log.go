package service

import (
	"gin-vue-admin/gin-vue-admin/dao"
	"gin-vue-admin/gin-vue-admin/response"
	"gin-vue-admin/gin-vue-admin/util/e"
	"github.com/gin-gonic/gin"
	"time"
)

type OperationLogService struct {
	ID              uint      `gorm:"type:varchar(20)" json:"id" form:"id"`
	Username        string    `gorm:"type:varchar(20);comment:'用户登录名'" json:"username"`
	Ip              string    `gorm:"type:varchar(20);comment:'Ip地址'" json:"ip"`
	IpLocation      string    `gorm:"type:varchar(20);comment:'Ip所在地'" json:"ipLocation"`
	Method          string    `gorm:"type:varchar(20);comment:'请求方式'" json:"method"`
	Path            string    `gorm:"type:varchar(100);comment:'访问路径'" json:"path"`
	Desc            string    `gorm:"type:varchar(100);comment:'说明'" json:"desc"`
	Status          int       `gorm:"type:int(4);comment:'响应状态码'" json:"status"`
	StartTime       time.Time `gorm:"type:datetime(3);comment:'发起时间'" json:"startTime"`
	TimeCost        int64     `gorm:"type:int(6);comment:'请求耗时(ms)'" json:"timeCost"`
	UserAgent       string    `gorm:"type:varchar(20);comment:'浏览器标识'" json:"userAgent"`
	PageNum         int       `json:"pageNum" form:"pageNum"`
	PageSize        int       `json:"pageSize" form:"pageSize"`
	OperationLogIds []uint    `json:"operationLogIds" form:"operationLogIds"`
}

// 获取操作日志列表
func (req OperationLogService) GetOperationLogs(c *gin.Context) response.Res {
	code := e.SUCCESS
	// 获取
	logs, total, err := dao.GetOperationLogs(req.Username, req.Ip, req.Path, req.Status, req.PageNum, req.PageSize)
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
			"total": total, "logs": logs,
		},
	}
}

// 批量删除操作日志
func (req OperationLogService) BatchDeleteOperationLogByIds(c *gin.Context) response.Res {
	code := e.SUCCESS
	// 删除接口
	err := dao.BatchDeleteOperationLogByIds(req.OperationLogIds)
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
