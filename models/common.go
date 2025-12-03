// Package models defines the data structures used in the MiniEye Intranet API.
package models

// Response represents a common API response structure.
type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// DataResponse represents a common API response with typed data.
type DataResponse struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// ListResponse represents a common API response for list data.
type ListResponse struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    ListData    `json:"data"`
}

// ListData represents list data structure.
type ListData struct {
	Total uint64      `json:"total"`
	List  interface{} `json:"list"`
}

// BaseDataInfo 基础带数据信息
type BaseDataInfo struct {
	// Error code | 错误代码
	Code int `json:"code"`

	// Message | 提示信息
	Msg string `json:"msg"`

	// Data | 数据
	Data string `json:"data,omitempty"`
}

// BaseListInfo 基础列表信息
type BaseListInfo struct {
	// The total number of data | 数据总数
	Total uint64 `json:"total"`

	// Data | 数据
	Data string `json:"data,omitempty"`
}

// BaseMsgResp 基础不带数据信息
type BaseMsgResp struct {
	// Error code | 错误代码
	Code int `json:"code"`

	// Message | 提示信息
	Msg string `json:"msg"`
}

// ErrorCode defines common error codes
const (
	ErrCodeSuccess       = 0    // 成功
	ErrCodeInvalidInput  = 4001 // 无效输入
	ErrCodeUnauthorized  = 4002 // 未授权
	ErrCodeForbidden     = 4003 // 禁止访问
	ErrCodeNotFound      = 4004 // 资源不存在
	ErrCodeServerError   = 5000 // 服务器错误
	ErrCodeDatabaseError = 5001 // 数据库错误
	ErrCodeNetworkError  = 5002 // 网络错误
)

// PageInfo 列表请求参数
type PageInfo struct {
	// Page number | 第几页
	Page uint64 `json:"page" validate:"required,number,gt=0"`

	// Page size | 单页数据行数
	PageSize uint64 `json:"pageSize" validate:"required,number,lt=100000"`
}

// UUIDReq 基础UUID参数请求
type UUIDReq struct {
	// ID
	// Required: true
	// Max length: 36
	Id string `json:"id" validate:"len=36"`
}

// UUIDsReq 基础UUID数组参数请求
type UUIDsReq struct {
	// Ids
	// Required: true
	Ids []string `json:"ids"`
}

// BaseUUIDInfo UUID基础信息
type BaseUUIDInfo struct {
	// UUID
	UUID *string `json:"uuid,optional"`

	// Create date | 创建日期
	CreatedAt *int64 `json:"createdAt,optional"`

	// Update date | 更新日期
	UpdatedAt *int64 `json:"updatedAt,optional"`
}