// Package models defines the data structures used in the MiniEye Intranet API.
package models

// UserInfo 用户信息
type UserInfo struct {
	// ID | 用户唯一标识（UUID）
	Id *string `json:"id,optional"`

	// Create date | 创建日期
	CreatedAt *int64 `json:"createdAt,optional"`

	// Update date | 更新日期
	UpdatedAt *int64 `json:"updatedAt,optional"`

	// Status | 状态
	Status *uint32 `json:"status,optional"`

	// Username | 用户名
	Username *string `json:"username,optional"`

	// Nickname | 昵称
	Nickname *string `json:"nickname,optional"`

	// Password | 密码
	Password *string `json:"password,optional"`

	// Description | 描述
	Description *string `json:"description,optional"`

	// HomePath | 首页
	HomePath *string `json:"homePath,optional"`

	// RoleIds | 角色ID列表
	RoleIds []uint64 `json:"roleIds,optional"`

	// Mobile | 手机号
	Mobile *string `json:"mobile,optional"`

	// Email | 邮箱
	Email *string `json:"email,optional"`

	// Avatar | 头像地址
	Avatar *string `json:"avatar,optional"`

	// DepartmentId | 部门ID
	DepartmentId *uint64 `json:"departmentId,optional"`

	// DepartmentName | 部门名称
	DepartmentName *string `json:"departmentName,optional"`

	// PositionIds | 职位ID列表
	PositionIds []uint64 `json:"positionId,optional"`

	// DingtalkId | 钉钉ID
	DingtalkId *string `json:"dingtalkId,optional"`

	// FeishuId | 飞书ID
	FeishuId *string `json:"feishuId,optional"`
}

// GetUsername 获取用户名
func (u *UserInfo) GetUsername() string {
	if u != nil && u.Username != nil && *u.Username != "" {
		return *u.Username
	}
	return ""
}

// GetNickname 获取昵称
func (u *UserInfo) GetNickname() string {
	if u != nil && u.Nickname != nil && *u.Nickname != "" {
		return *u.Nickname
	}
	return ""
}

// UserListReq 用户列表请求
type UserListReq struct {
	// PageInfo
	PageInfo

	// User Name | 用户名
	Username *string `json:"username,optional"`

	// User's nickname | 用户的昵称
	Nickname *string `json:"nickname,optional"`

	// User's mobile phone number | 用户的手机号码
	Mobile *string `json:"mobile,optional"`

	// The user's email address | 用户的邮箱
	Email *string `json:"email,optional"`

	// User's role ID | 用户的角色ID
	RoleIds []uint64 `json:"roleIds,optional"`

	// The user's department ID | 用户所属部门ID
	DepartmentId *uint64 `json:"departmentId,optional"`

	// User's position id | 用户的职位ID
	PositionId *uint64 `json:"positionId,optional"`

	// Description | 描述
	Description *string `json:"description,optional"`
}

// UserListRsp 用户列表响应
type UserListRsp struct {
	// The total number of data | 数据总数
	Total uint64 `json:"total"`

	// User list | 用户列表
	Data []UserInfo `json:"data"`
}

// UserUpdateReq 更新用户请求
type UserUpdateReq struct {
	// UUID
	// Required: true
	UUID string `json:"uuid" validate:"required"`

	// Real name | 真实姓名
	RealName string `json:"realName,optional"`

	// Email | 邮箱
	Email string `json:"email,optional"`

	// Phone | 手机号
	Phone string `json:"phone,optional"`

	// Department | 部门
	Department string `json:"department,optional"`

	// Role | 角色
	Role string `json:"role,optional"`

	// IsEnable | 是否启用
	IsEnable bool `json:"isEnable,optional"`
}

// UserCreateReq 创建用户请求
type UserCreateReq struct {
	// User name | 用户名
	// Required: true
	Username string `json:"username" validate:"required,min=3,max=50"`

	// Password | 密码
	// Required: true
	Password string `json:"password" validate:"required,min=6"`

	// Real name | 真实姓名
	// Required: true
	RealName string `json:"realName" validate:"required,min=1,max=50"`

	// Email | 邮箱
	Email string `json:"email,optional" validate:"omitempty,email"`

	// Phone | 手机号
	Phone string `json:"phone,optional" validate:"omitempty,len=11"`

	// Department | 部门
	Department string `json:"department,optional"`

	// Position | 职位
	Position string `json:"position,optional"`

	// Role | 角色
	Role string `json:"role,optional"`

	// IsEnable | 是否启用
	IsEnable bool `json:"isEnable,optional"`
}

// UserDelReq 删除用户请求
type UserDelReq struct {
	// UUID列表
	// Required: true
	UUIDs []string `json:"uuids" validate:"required"`
}

// UserResetPwdReq 重置密码请求
type UserResetPwdReq struct {
	// UUID
	// Required: true
	UUID string `json:"uuid" validate:"required"`

	// Password | 密码
	// Required: true
	Password string `json:"password" validate:"required,min=6"`
}

// ChangePasswordReq 修改密码请求
type ChangePasswordReq struct {
	// Old password | 旧密码
	// Required: true
	OldPassword string `json:"oldPassword" validate:"required"`

	// New password | 新密码
	// Required: true
	NewPassword string `json:"newPassword" validate:"required,min=6"`
}

// UpdateProfileReq 更新个人资料请求
type UpdateProfileReq struct {
	// Real name | 真实姓名
	RealName string `json:"realName,optional" validate:"omitempty,min=1,max=50"`

	// Email | 邮箱
	Email string `json:"email,optional" validate:"omitempty,email"`

	// Phone | 手机号
	Phone string `json:"phone,optional" validate:"omitempty,len=11"`

	// Avatar | 头像
	Avatar string `json:"avatar,optional"`

	// Position | 职位
	Position string `json:"position,optional"`
}
