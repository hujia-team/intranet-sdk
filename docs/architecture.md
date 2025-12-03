# MiniEye Intranet SDK 架构设计

## 整体架构

```
├── intranet.go          # SDK主入口，提供客户端初始化和核心功能
├── client/              # HTTP客户端实现
│   └── http_client.go   # 处理HTTP通信和认证
├── models/              # 数据模型定义
│   ├── user.go          # 用户相关模型
│   └── common.go        # 通用数据结构
├── services/            # 业务服务层
│   ├── user_service.go  # 用户服务
│   └── connector_service.go # 连接器服务
├── utils/               # 工具函数
│   ├── errors.go        # 错误处理
│   ├── logger.go        # 日志工具
│   └── sts.go           # STS认证相关工具
├── examples/            # 使用示例
│   ├── user_example.go  # 用户信息示例
│   └── connector_example.go # 连接器示例
└── go.mod               # Go模块定义
```

## API接口设计

### 1. 核心接口

#### 客户端初始化
```go
// NewClient 创建一个新的MiniEye Intranet API客户端
func NewClient(options ...Option) (*Client, error)
```

#### 配置选项
```go
// WithBaseURL 设置API基础URL
func WithBaseURL(url string) Option

// WithAPIKey 设置API密钥
func WithAPIKey(apiKey string) Option

// WithUserAgent 设置用户代理
func WithUserAgent(userAgent string) Option

// WithAccessKeyID 设置STS认证的访问密钥ID
func WithAccessKeyID(accessKeyID string) Option

// WithAccessKeySecret 设置STS认证的访问密钥密钥
func WithAccessKeySecret(accessKeySecret string) Option
```

### 2. 用户服务

```go
// UserService 提供用户相关的功能
type UserService interface {
    // GetUserInfo 获取当前用户信息
    GetUserInfo() (*models.UserInfo, error)
}
```

### 3. 连接器服务

```go
// ConnectorService 提供连接器相关的功能
type ConnectorService interface {
    // SendKafkaMessage sends a message to Kafka.
    // topic: the Kafka topic to send the message to
    // message: the message to send, can be any type that can be marshaled to JSON
    SendKafkaMessage(topic string, message any) (models.BaseMsgResp, error)

}
```

## 数据模型设计

### 通用模型
```go
// Response 通用响应结构
type Response struct {
    Code    int         `json:"code"`
    Message string      `json:"message"`
    Data    interface{} `json:"data,omitempty"`
}

// DataResponse 带类型数据的通用API响应
type DataResponse struct {
    Code    int         `json:"code"`
    Message string      `json:"message"`
    Data    interface{} `json:"data,omitempty"`
}

// ListResponse 列表数据的通用API响应
type ListResponse struct {
    Code    int         `json:"code"`
    Message string      `json:"message"`
    Data    ListData    `json:"data"`
}

// ListData 列表数据结构
type ListData struct {
    Total uint64      `json:"total"`
    List  interface{} `json:"list"`
}

// BaseMsgResp 基础不带数据信息
type BaseMsgResp struct {
    Code int    `json:"code"`
    Msg  string `json:"msg"`
}

// PageInfo 列表请求参数
type PageInfo struct {
    Page     uint64 `json:"page" validate:"required,number,gt=0"`
    PageSize uint64 `json:"pageSize" validate:"required,number,lt=100000"`
}
```

### 用户模型
```go
// UserInfo 用户信息
type UserInfo struct {
    // UUID | 用户唯一标识
    UserID string `json:"userId,optional"`

    // User name | 用户名
    Username string `json:"username,optional"`

    // Nickname | 昵称
    Nickname string `json:"nickname,optional"`

    // Avatar | 头像
    Avatar string `json:"avatar,optional"`

    // HomePath | 主目录路径
    HomePath string `json:"homePath,optional"`

    // RoleName | 角色名称
    RoleName string `json:"roleName,optional"`

    // DepartmentName | 部门名称
    DepartmentName string `json:"departmentName,optional"`
}
```

## 错误处理设计

SDK使用SDKError作为自定义错误类型，提供详细的错误信息和错误码。错误类型定义如下：

```go
// SDKError SDK错误类型
type SDKError struct {
    Code    ErrorCode
    Message string
    Err     error
}
```

SDK提供了多种便捷的错误创建函数，如NewInvalidInputError、NewUnauthorizedError等，具体实现请参考utils/errors.go文件。

## 认证机制

SDK支持多种认证方式，优先级从高到低为：
1. 访问令牌认证（Bearer Token）：通过登录获取的authToken
2. STS认证：使用AccessKeyID和AccessKeySecret
3. API密钥认证：直接使用API密钥作为Bearer Token

## 日志和调试

SDK将提供可配置的日志级别，方便调试和问题排查。