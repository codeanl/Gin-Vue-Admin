package service

import (
	"gin-vue-admin/gin-vue-admin/dao"
	"gin-vue-admin/gin-vue-admin/model"
	"gin-vue-admin/gin-vue-admin/response"
	"gin-vue-admin/gin-vue-admin/util"
	"gin-vue-admin/gin-vue-admin/util/e"
	"github.com/gin-gonic/gin"
	"github.com/thoas/go-funk"
)

type UserService struct {
	ID           uint   `gorm:"type:varchar(20)" json:"id" form:"id"`
	Username     string `gorm:"type:varchar(20);not null;unique" json:"username" form:"username"`
	Password     string `gorm:"size:255;not null" json:"password" form:"password"`
	Mobile       string `gorm:"type:varchar(11);not null;unique" json:"mobile" form:"mobile"`
	Avatar       string `gorm:"type:varchar(255)" json:"avatar"`
	Nickname     string `gorm:"type:varchar(20)" json:"nickname" form:"nickname"`
	Introduction string `gorm:"type:varchar(255)" json:"introduction"`
	Status       uint   `gorm:"type:tinyint(1);default:1;comment:'1正常, 2禁用'" json:"status"`
	Creator      string `gorm:"type:varchar(20);" json:"creator"`
	PageNum      uint   `json:"pageNum" form:"pageNum"`
	PageSize     uint   `json:"pageSize" form:"pageSize"`
	RoleIds      []uint `form:"roleIds" json:"roleIds" form:"roleIds"`
	UserIds      []uint `json:"userIds" form:"userIds"`
}

// Login 用户登陆
func (req UserService) Login(c *gin.Context) response.Res {
	code := e.SUCCESS
	user, exist, err := dao.ExistOrNotByUserName(req.Username)
	if !exist { //如果查询不到，返回相应的错误
		code = e.UserNotExist
		return response.Res{
			Code: code,
			Msg:  e.GetMsg(code),
		}
	}
	if user.CheckPassword(req.Password) == false {
		code = e.PasswordError
		return response.Res{
			Code: code,
			Msg:  e.GetMsg(code),
		}
	}
	token, err := util.GenerateToken(user.ID, req.Username)
	if err != nil {
		code = e.ERROR
		return response.Res{
			Code: code,
			Msg:  e.GetMsg(code),
		}
	}
	return response.Res{
		Code: code,
		Data: response.TokenData{User: user, Authorization: token},
		Msg:  e.GetMsg(code),
	}
}

// Logout 退出登录
func (req UserService) Logout(c *gin.Context, Authorization string) response.Res {
	code := e.SUCCESS
	util.InvalidToken(Authorization)
	return response.Res{
		Code: code,
		Msg:  e.GetMsg(code),
	}
}

// GetMyInfo 获取自己的详细信息
func (req UserService) GetMyInfo(c *gin.Context, myid uint) response.Res {
	code := e.SUCCESS
	data, err := dao.GetUserById(myid)
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
		Data: data,
	}
}

// GetUsers 获取用户列表
func (req UserService) GetUsers(c *gin.Context) response.Res {
	code := e.SUCCESS
	// 获取
	users, total, err := dao.GetUsers(req.Username, req.Nickname, req.Mobile, req.Status, int(req.PageNum), int(req.PageSize))
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
			"users": response.ToUsersDto(users),
			"total": total,
		},
	}
}

// 更新用户登录密码
func (req *UserService) ChangePwd(c *gin.Context, myid uint, oldpassword string, newpassword string) response.Res {
	code := e.SUCCESS
	if oldpassword == "" || newpassword == "" {
		code = e.InvalidParams
		return response.Res{
			Code: code,
			Msg:  e.GetMsg(code),
		}
	}
	user, _ := dao.GetUserById(myid)
	if user.CheckPassword(oldpassword) == false {
		code = e.PasswordError
		return response.Res{
			Code: code,
			Msg:  e.GetMsg(code),
		}
	} else {
		user.SetPassword(newpassword)
		err := dao.UpdateUserById(myid, user)
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
}

// 创建用户
func (req *UserService) CreateUser(c *gin.Context, myid uint) response.Res {
	code := e.SUCCESS
	// 当前用户角色排序最小值（最高等级角色）以及当前用户
	currentRoleSortMin, ctxUser, err := dao.GetCurrentUserMinRoleSort(c, myid)
	if err != nil {
		code = e.ERROR
		return response.Res{
			Code:  code,
			Msg:   e.GetMsg(code),
			Error: err.Error(),
		}
	}
	// 根据角色id获取角色
	roles, err := dao.GetRolesByIds(req.RoleIds)
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
			Code: code,
			Msg:  e.GetMsg(code),
		}
	}
	var reqRoleSorts []int
	for _, role := range roles {
		reqRoleSorts = append(reqRoleSorts, int(role.Sort))
	}
	// 前端传来用户角色排序最小值（最高等级角色）
	reqRoleSortMin := uint(funk.MinInt(reqRoleSorts))

	// 当前用户的角色排序最小值 需要小于 前端传来的角色排序最小值（用户不能创建比自己等级高的或者相同等级的用户）
	if currentRoleSortMin >= reqRoleSortMin {
		code = e.CreatUserRoleGreater
		return response.Res{
			Code: code,
			Msg:  e.GetMsg(code),
		}
	}
	user := model.User{
		Username:     req.Username,
		Password:     req.Password,
		Mobile:       req.Mobile,
		Avatar:       req.Avatar,
		Nickname:     req.Nickname,
		Introduction: req.Introduction,
		Status:       req.Status,
		Creator:      ctxUser.Username,
		Roles:        roles,
	}
	// 密码为空就默认123456
	if req.Password == "" {
		req.Password = "123456"
	}
	// 密码不为空就解密
	if req.Password != "" {
		if len(req.Password) < 6 {
			code = e.PasswordParams6
			return response.Res{
				Code:  code,
				Msg:   e.GetMsg(code),
				Error: "密码长度至少为6位",
			}
		}
		err := user.SetPassword(req.Password)
		if err != nil {
			code = e.ERROR
			return response.Res{
				Code:  code,
				Msg:   e.GetMsg(code),
				Error: err.Error(),
			}
		}
	}
	err = dao.CreateUser(&user)
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
		Data: user,
	}
}

// 更新用户
func (req *UserService) UpdateUserById(c *gin.Context, myid uint) response.Res {
	code := e.SUCCESS
	//获取path中的userId
	if req.ID <= 0 {
		code = e.UserIDError
		return response.Res{
			Code: code,
			Msg:  e.GetMsg(code),
		}
	}
	// 根据path中的userId获取用户信息
	oldUser, err := dao.GetUserById(req.ID)
	if err != nil {
		code = e.ERROR
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
			Code:  code,
			Msg:   e.GetMsg(code),
			Error: err.Error(),
		}
	}
	// 获取当前用户的所有角色
	currentRoles := ctxUser.Roles
	// 获取当前用户角色的排序，和前端传来的角色排序做比较
	var currentRoleSorts []int
	// 当前用户角色ID集合
	var currentRoleIds []uint
	for _, role := range currentRoles {
		currentRoleSorts = append(currentRoleSorts, int(role.Sort))
		currentRoleIds = append(currentRoleIds, role.ID)
	}
	// 当前用户角色排序最小值（最高等级角色）
	currentRoleSortMin := funk.MinInt(currentRoleSorts)
	// 获取前端传来的用户角色id
	reqRoleIds := req.RoleIds
	// 根据角色id获取角色
	roles, err := dao.GetRolesByIds(reqRoleIds)
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
	var reqRoleSorts []int
	for _, role := range roles {
		reqRoleSorts = append(reqRoleSorts, int(role.Sort))
	}
	// 前端传来用户角色排序最小值（最高等级角色）
	reqRoleSortMin := funk.MinInt(reqRoleSorts)
	user := model.User{
		Model:        oldUser.Model,
		Username:     req.Username,
		Password:     oldUser.Password,
		Mobile:       req.Mobile,
		Avatar:       req.Avatar,
		Nickname:     req.Nickname,
		Introduction: req.Introduction,
		Status:       req.Status,
		Creator:      ctxUser.Username,
		Roles:        roles,
	}
	// 判断是更新自己还是更新别人
	if req.ID == ctxUser.ID {
		// 如果是更新自己
		// 不能禁用自己
		if req.Status == 2 {
			code = e.StatusNotMe
			return response.Res{
				Code: code,
				Msg:  e.GetMsg(code),
			}
		}
		// 不能更改自己的角色
		reqDiff, currentDiff := funk.Difference(req.RoleIds, currentRoleIds)
		if len(reqDiff.([]uint)) > 0 || len(currentDiff.([]uint)) > 0 {
			code = e.RoleNotMe
			return response.Res{
				Code: code,
				Msg:  e.GetMsg(code),
			}
		}
		// 不能更新自己的密码，只能在个人中心更新
		if req.Password != "" {
			code = e.PasswordTo
			return response.Res{
				Code: code,
				Msg:  e.GetMsg(code),
			}
		}
		// 密码赋值
		user.Password = ctxUser.Password
	} else {
		// 如果是更新别人
		// 用户不能更新比自己角色等级高的或者相同等级的用户
		// 根据path中的userIdID获取用户角色排序最小值
		minRoleSorts, err := dao.GetUserMinRoleSortsByIds([]uint{req.ID})
		if err != nil || len(minRoleSorts) == 0 {
			code = e.ERROR
			return response.Res{
				Code: code,
				Msg:  e.GetMsg(code),
			}
		}
		if currentRoleSortMin >= minRoleSorts[0] {
			code = e.CreatUserRoleGreater
			return response.Res{
				Code: code,
				Msg:  e.GetMsg(code),
			}
		}
		// 用户不能把别的用户角色等级更新得比自己高或相等
		if currentRoleSortMin >= reqRoleSortMin {
			code = e.RoleGreaterMe
			return response.Res{
				Code: code,
				Msg:  e.GetMsg(code),
			}
		}
		// 密码赋值
		if req.Password != "" {
			if err = user.SetPassword(req.Password); err != nil {
				code = e.ERROR
				return response.Res{
					Code: code,
					Msg:  e.GetMsg(code),
				}
			}
		}
	}
	// 更新用户
	err = dao.UpdateUser(&user)
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

// 批量删除用户
func (req *UserService) BatchDeleteUserByIds(c *gin.Context, myid uint) response.Res {
	code := e.SUCCESS

	// 前端传来的用户ID
	reqUserIds := req.UserIds
	// 根据用户ID获取用户角色排序最小值
	roleMinSortList, err := dao.GetUserMinRoleSortsByIds(reqUserIds)
	if err != nil || len(roleMinSortList) == 0 {
		code = e.ERROR
		return response.Res{
			Code:  code,
			Msg:   e.GetMsg(code),
			Error: "根据用户ID获取用户角色排序最小值失败",
		}
	}

	// 当前用户角色排序最小值（最高等级角色）以及当前用户
	minSort, ctxUser, err := dao.GetCurrentUserMinRoleSort(c, myid)
	if err != nil {
		code = e.ERROR
		return response.Res{
			Code:  code,
			Msg:   e.GetMsg(code),
			Error: err.Error(),
		}
	}
	currentRoleSortMin := int(minSort)

	// 不能删除自己
	if funk.Contains(reqUserIds, ctxUser.ID) {
		code = e.ERROR
		return response.Res{
			Code:  code,
			Msg:   e.GetMsg(code),
			Error: "用户不能删除自己",
		}
	}

	// 不能删除比自己角色排序低(等级高)的用户
	for _, sort := range roleMinSortList {
		if currentRoleSortMin >= sort {
			code = e.ERROR
			return response.Res{
				Code:  code,
				Msg:   e.GetMsg(code),
				Error: "用户不能删除比自己角色等级高的用户",
			}
		}
	}

	err = dao.BatchDeleteUserByIds(reqUserIds)
	if err != nil {
		code = e.ERROR
		return response.Res{
			Code:  code,
			Msg:   e.GetMsg(code),
			Error: "删除用户失败: " + err.Error(),
		}
	}

	return response.Res{
		Code: code,
		Msg:  e.GetMsg(code),
	}

}
