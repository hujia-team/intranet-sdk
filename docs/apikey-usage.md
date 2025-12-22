# ApiKey 管理功能使用指南

本文档介绍如何使用 intranet-sdk 的 ApiKey 管理功能。

## 功能概述

ApiKey 服务提供了以下功能:

- 获取 ApiKey 列表
- 根据 ID 获取单个 ApiKey 详情
- 创建新的 ApiKey
- 更新现有的 ApiKey
- 删除 ApiKey

## 快速开始

### 1. 初始化客户端

```go
import "github.com/hujia-team/intranet-sdk"

// 使用 STS 认证
client, err := intranet.NewClient(
    intranet.WithAccessKeyID("your_access_key_id"),
    intranet.WithAccessKeySecret("your_access_key_secret"),
)
if err != nil {
    // 处理错误
}
```

### 2. 获取 ApiKey 列表

```go
import "github.com/hujia-team/intranet-sdk/models"

// 创建请求参数
req := &models.ApiKeyListReq{
    PageInfo: models.PageInfo{
        Page:     1,
        PageSize: 10,
    },
    // 可选的过滤条件
    Domain: "your-domain",  // 按域筛选
}

// 获取列表
resp, err := client.ApiKey.GetApiKeyList(req)
if err != nil {
    // 处理错误
}

// 处理结果
fmt.Printf("总数: %d\n", resp.Total)
for _, apiKey := range resp.List {
    fmt.Printf("ID: %d, Name: %s\n", apiKey.ID, apiKey.Name)
}
```

### 3. 获取单个 ApiKey 详情

```go
apiKey, err := client.ApiKey.GetApiKeyByID(123)
if err != nil {
    // 处理错误
}

fmt.Printf("ApiKey: %s\n", apiKey.Name)
fmt.Printf("Domain: %s\n", apiKey.Domain)
fmt.Printf("BaseURL: %s\n", apiKey.BaseURL)

// 查看统计信息
if apiKey.Stats != nil && apiKey.Stats.Usage != nil {
    fmt.Printf("请求数: %d\n", apiKey.Stats.Usage.Requests)
    fmt.Printf("总 Token 数: %d\n", apiKey.Stats.Usage.AllTokens)
    fmt.Printf("成本: %.4f\n", apiKey.Stats.Usage.Cost)
}
```

### 4. 创建新的 ApiKey

```go
newApiKey := &models.ApiKeyInfo{
    Name:        "我的 ApiKey",
    Description: "用于测试的 ApiKey",
    Domain:      "your-domain",
    BaseURL:     "https://api.example.com",
}

id, err := client.ApiKey.CreateApiKey(newApiKey)
if err != nil {
    // 处理错误
}

fmt.Printf("创建成功,新 ApiKey ID: %d\n", id)
```

### 5. 更新 ApiKey

```go
updateApiKey := &models.ApiKeyInfo{
    ID:          123,
    Name:        "更新后的名称",
    Description: "更新后的描述",
}

err := client.ApiKey.UpdateApiKey(updateApiKey)
if err != nil {
    // 处理错误
}

fmt.Println("更新成功!")
```

### 6. 删除 ApiKey

```go
// 删除单个 ApiKey
err := client.ApiKey.DeleteApiKey([]uint64{123})
if err != nil {
    // 处理错误
}

// 批量删除多个 ApiKey
err = client.ApiKey.DeleteApiKey([]uint64{123, 456, 789})
if err != nil {
    // 处理错误
}

fmt.Println("删除成功!")
```

## 数据结构说明

### ApiKeyInfo

```go
type ApiKeyInfo struct {
    ID          uint64        // ApiKey ID
    CreatedAt   int64         // 创建时间
    UpdatedAt   int64         // 更新时间
    Name        string        // ApiKey 名称
    Description string        // ApiKey 描述
    BaseURL     string        // 基础 URL
    Token       string        // ApiKey 令牌
    Domain      string        // 域
    Sub         string        // 主体
    SubType     string        // 主体类型
    IsOwner     bool          // 是否是所有者
    IsAdmin     bool          // 是否是管理员
    HasRead     bool          // 是否有读权限
    HasWrite    bool          // 是否有写权限
    HasUse      bool          // 是否有使用权限
    Stats       *ApiKeyStats  // 统计信息
}
```

### ApiKeyStats

```go
type ApiKeyStats struct {
    APIID        string     // API ID
    Name         string     // 名称
    IsActive     bool       // 是否活跃
    Usage        *UsageData // 总使用统计
    DailyUsage   *UsageData // 每日使用统计
    MonthlyUsage *UsageData // 每月使用统计
}
```

### UsageData

```go
type UsageData struct {
    Requests          int64   // 请求数
    InputTokens       int64   // 输入 token 数
    OutputTokens      int64   // 输出 token 数
    CacheCreateTokens int64   // 缓存创建 token 数
    CacheReadTokens   int64   // 缓存读取 token 数
    AllTokens         int64   // 总 token 数
    Cost              float64 // 成本
    FormattedCost     string  // 格式化后的成本
}
```

## 完整示例

参考 [examples/apikey_example.go](../examples/apikey_example.go) 获取完整的使用示例。

## 权限说明

使用 ApiKey 管理功能需要相应的权限:

- **读取权限** (HasRead): 可以查看 ApiKey 列表和详情
- **写入权限** (HasWrite): 可以创建和更新 ApiKey
- **删除权限** (IsAdmin): 可以删除 ApiKey
- **所有者权限** (IsOwner): 拥有所有权限

## 错误处理

所有的 ApiKey 操作都可能返回错误,建议进行适当的错误处理:

```go
resp, err := client.ApiKey.GetApiKeyList(req)
if err != nil {
    if sdkErr, ok := err.(*utils.SDKError); ok {
        log.Printf("SDK 错误: %s", sdkErr.Error())
    } else {
        log.Printf("其他错误: %v", err)
    }
    return
}
```

## 注意事项

1. ApiKey Token 是敏感信息,请妥善保管
2. 删除操作不可逆,请谨慎操作
3. 统计信息 (Stats) 可能为空,使用前请检查
4. 部分字段可能是可选的,根据实际业务场景使用
