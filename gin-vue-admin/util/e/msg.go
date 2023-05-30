package e

var MsgFlags = map[int]string{
	SUCCESS:       "成功",
	ERROR:         "错误",
	InvalidParams: "输入的参数不全",
	//user
	UserNotExist:         "用户不存在",
	PasswordError:        "密码错误",
	UserRoleError:        "获取角色信息失败",
	CreatUserRoleGreater: "用户等级比自己等级高的或者相同",
	PasswordParams6:      "密码长度不足6位",
	UserIDError:          "用户ID不正确",
	StatusNotMe:          "不能禁用自己",
	RoleNotMe:            "不能更改自己的角色",
	PasswordTo:           "请到个人中心更新自身密码",
	RoleGreaterMe:        "用户角色等级更新得比自己高或相等",
	//role
	RoleNameExist:    "角色已存在",
	CreatRoleGreater: "创建比自己等级高或相同等级的角色",
	RoleIDError:      "角色ID不正确",
	//api
	ApiIDError: "角色ID不正确",
}

// GetMsg 获取状态码对应信息
func GetMsg(code int) string {
	msg, ok := MsgFlags[code]
	if ok {
		return msg
	}
	return MsgFlags[ERROR]
}
