package service

import (
	"gin-vue-admin/gin-vue-admin/dao"
	"gin-vue-admin/gin-vue-admin/model"
	"gin-vue-admin/gin-vue-admin/response"
	"gin-vue-admin/gin-vue-admin/util/e"
	"github.com/gin-gonic/gin"
)

type ApiService struct {
	ID       uint   `gorm:"type:varchar(20)" json:"id" form:"id"`
	Method   string `gorm:"type:varchar(20);comment:'请求方式'" json:"method" form:"method"`
	Path     string `gorm:"type:varchar(100);comment:'访问路径'" json:"path" form:"path"`
	Category string `gorm:"type:varchar(50);comment:'所属类别'" json:"category"  form:"category"`
	Desc     string `gorm:"type:varchar(100);comment:'说明'" json:"desc" form:"desc"`
	Creator  string `gorm:"type:varchar(20);comment:'创建人'" json:"creator" form:"creator"`
	PageNum  int    `json:"pageNum" form:"pageNum"`
	PageSize int    `json:"pageSize" form:"pageSize"`
	ApiIds   []uint `json:"apiIds" form:"apiIds"`
}

// GetApis 获取接口列表
func (req ApiService) GetApis(c *gin.Context) response.Res {
	code := e.SUCCESS
	// 获取
	apis, total, err := dao.GetApis(req.Method, req.Path, req.Category, req.Creator, req.PageNum, req.PageSize)
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
			"apis":  apis,
		},
	}
}

// GetApiTree 获取接口树(按接口Category字段分类)
func (req ApiService) GetApiTree(c *gin.Context) response.Res {
	code := e.SUCCESS
	tree, err := dao.GetApiTree()
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
		Data: tree,
	}
}

// CreateApi 创建接口
func (req ApiService) CreateApi(c *gin.Context, myid uint) response.Res {
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
	api := model.Api{
		Method:   req.Method,
		Path:     req.Path,
		Category: req.Category,
		Desc:     req.Desc,
		Creator:  ctxUser.Username,
	}
	// 创建接口
	err = dao.CreateApi(&api)
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
		Data: api,
	}
}

// UpdateApiById 更新接口
func (req ApiService) UpdateApiById(c *gin.Context, myid uint) response.Res {
	code := e.SUCCESS
	// 获取路径中的apiId
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
	api := model.Api{
		Method:   req.Method,
		Path:     req.Path,
		Category: req.Category,
		Desc:     req.Desc,
		Creator:  ctxUser.Username,
	}
	err = dao.UpdateApiById(req.ID, &api)
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
		Data: api,
	}
}

// BatchDeleteApiByIds 批量删除接口
func (req ApiService) BatchDeleteApiByIds(c *gin.Context, myid uint) response.Res {
	code := e.SUCCESS
	// 删除接口
	err := dao.BatchDeleteApiByIds(req.ApiIds)
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
