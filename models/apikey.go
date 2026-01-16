// Package models defines the data structures used in the MiniEye Intranet API.
package models

// ApiKeyInfo ApiKey 信息
type ApiKeyInfo struct {
	// ID | ApiKey唯一标识
	ID uint64 `json:"id,optional"`

	// CreatedAt | 创建时间
	CreatedAt int64 `json:"createdAt,optional"`

	// UpdatedAt | 更新时间
	UpdatedAt int64 `json:"updatedAt,optional"`

	// Name | ApiKey名称
	Name string `json:"name,optional"`

	// Description | ApiKey描述
	Description string `json:"description,optional"`

	// BaseURL | ApiKey基础URL
	BaseURL string `json:"baseUrl,optional"`

	// Token | ApiKey令牌
	Token string `json:"token,optional"`

	// Domain | 域
	Domain string `json:"domain,optional"`

	// Sub | 主体
	Sub string `json:"sub,optional"`

	// SubType | 主体类型
	SubType string `json:"subType,optional"`

	// IsOwner | 是否是所有者
	IsOwner bool `json:"isOwner,optional"`

	// IsAdmin | 是否是管理员
	IsAdmin bool `json:"isAdmin,optional"`

	// HasRead | 是否有读权限
	HasRead bool `json:"hasRead,optional"`

	// HasWrite | 是否有写权限
	HasWrite bool `json:"hasWrite,optional"`

	// HasUse | 是否有使用权限
	HasUse bool `json:"hasUse,optional"`

	// Stats | ApiKey统计信息
	Stats *ApiKeyStats `json:"stats,optional"`
}

// ApiKeyStats ApiKey 统计信息
type ApiKeyStats struct {
	// APIID | API ID
	APIID string `json:"apiId,optional"`

	// Name | 名称
	Name string `json:"name,optional"`

	// IsActive | 是否活跃
	IsActive bool `json:"isActive,optional"`

	// Usage | 使用统计
	Usage *UsageData `json:"usage,optional"`

	// DailyUsage | 每日使用统计
	DailyUsage *UsageData `json:"dailyUsage,optional"`

	// MonthlyUsage | 每月使用统计
	MonthlyUsage *UsageData `json:"monthlyUsage,optional"`
}

// UsageData 使用统计数据
type UsageData struct {
	// Requests | 请求数
	Requests int64 `json:"requests,optional"`

	// InputTokens | 输入token数
	InputTokens int64 `json:"inputTokens,optional"`

	// OutputTokens | 输出token数
	OutputTokens int64 `json:"outputTokens,optional"`

	// CacheCreateTokens | 缓存创建token数
	CacheCreateTokens int64 `json:"cacheCreateTokens,optional"`

	// CacheReadTokens | 缓存读取token数
	CacheReadTokens int64 `json:"cacheReadTokens,optional"`

	// AllTokens | 总token数
	AllTokens int64 `json:"allTokens,optional"`

	// Cost | 成本
	Cost float64 `json:"cost,optional"`

	// FormattedCost | 格式化后的成本
	FormattedCost string `json:"formattedCost,optional"`
}

// ApiKeyListReq ApiKey列表请求参数
type ApiKeyListReq struct {
	// PageInfo
	PageInfo

	// CreatedAt | 创建时间
	CreatedAt int64 `json:"createdAt,optional"`

	// UpdatedAt | 更新时间
	UpdatedAt int64 `json:"updatedAt,optional"`

	// Name | 名称
	Name string `json:"name,optional"`

	// Description | 描述
	Description string `json:"description,optional"`

	// BaseURL | 基础URL
	BaseURL string `json:"baseUrl,optional"`

	// Token | 令牌
	Token string `json:"token,optional"`

	// Domain | 域
	Domain string `json:"domain,optional"`

	// Sub | 主体
	Sub string `json:"sub,optional"`

	// SubType | 主体类型
	SubType string `json:"subType,optional"`
}

// ApiKeyListResp ApiKey列表响应
type ApiKeyListResp struct {
	// Total | 数据总数
	Total uint64 `json:"total"`

	// List | ApiKey列表
	List []ApiKeyInfo `json:"list"`
}

// Sub2ApiGroupInfo Sub2Api 分组信息
type Sub2ApiGroupInfo struct {
	// ID | 分组ID
	ID int64 `json:"id"`

	// Name | 分组名称
	Name string `json:"name"`

	// Description | 分组描述
	Description string `json:"description"`

	// CreatedAt | 创建时间
	CreatedAt string `json:"createdAt"`

	// UpdatedAt | 更新时间
	UpdatedAt string `json:"updatedAt"`

	// Domain | 域
	Domain string `json:"domain"`

	// Sub | 主体
	Sub string `json:"sub"`

	// SubType | 主体类型
	SubType string `json:"subType"`

	// IsOwner | 是否是所有者
	IsOwner bool `json:"isOwner,optional"`

	// IsAdmin | 是否是管理员
	IsAdmin bool `json:"isAdmin,optional"`

	// HasRead | 是否有读权限
	HasRead bool `json:"hasRead,optional"`

	// HasWrite | 是否有写权限
	HasWrite bool `json:"hasWrite,optional"`

	// HasUse | 是否有使用权限
	HasUse bool `json:"hasUse,optional"`

	// DailyLimit | 日配额限制（单位：分，即美分）
	DailyLimit int64 `json:"dailyLimit,optional"`

	// DailyUsed | 日配额已使用（单位：分）
	DailyUsed int64 `json:"dailyUsed,optional"`

	// WeeklyLimit | 周配额限制（单位：分）
	WeeklyLimit int64 `json:"weeklyLimit,optional"`

	// WeeklyUsed | 周配额已使用（单位：分）
	WeeklyUsed int64 `json:"weeklyUsed,optional"`
}

// GetAvailableGroupsResp 获取可用分组响应
type GetAvailableGroupsResp struct {
	// Data | 分组列表
	Data []Sub2ApiGroupInfo `json:"data"`
}

// GetCurrentGroupResp 获取当前分组响应
type CurrentGroupResp struct {
	// BaseDataInfo | 基础数据信息
	BaseDataInfo

	// Data | 当前分组（可能为空）
	Data *Sub2ApiGroupInfo `json:"data,optional"`
}

// SwitchGroupReq 切换分组请求
type SwitchGroupReq struct {
	// GroupID | 目标分组 ID
	GroupID int64 `json:"groupId"`
}
