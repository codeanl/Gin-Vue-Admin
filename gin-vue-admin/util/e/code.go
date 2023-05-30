package e

const (
	SUCCESS       = 200
	InvalidParams = 201
	ERROR         = 400
	//user
	UserNotExist         = 10001
	PasswordError        = 10002
	UserRoleError        = 10003
	CreatUserRoleGreater = 10004
	PasswordParams6      = 10005
	UserIDError          = 10006
	StatusNotMe          = 10007
	RoleNotMe            = 10008
	PasswordTo           = 10009
	RoleGreaterMe        = 10010
	//role
	RoleNameExist    = 20001
	CreatRoleGreater = 20002
	RoleIDError      = 20003
	//api
	ApiIDError = 30001
)
