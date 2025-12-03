# 内网SDK使用指南

本指南提供了内网SDK的详细使用说明，包括初始化、认证、API调用和错误处理等方面。

## 1. 环境准备

### 1.1 安装SDK

使用Go模块安装SDK：

```bash
go get github.com/hujia-team/intranet-sdk
```

### 1.2 配置环境变量

SDK支持通过环境变量配置认证信息：

```bash
# STS认证相关环境变量
export INTRANET_ACCESS_KEY_ID=your_access_key_id
export INTRANET_ACCESS_KEY_SECRET=your_access_key_secret

# 或者使用以下环境变量名称
export STS_ACCESS_KEY_ID=your_access_key_id
export STS_ACCESS_KEY_SECRET=your_access_key_secret
export STS_SECURITY_TOKEN=your_security_token
```

## 2. 初始化客户端

SDK提供了多种初始化方式，适应不同的使用场景。

### 2.1 基本初始化（使用环境变量）

```go
import "github.com/hujia-team/intranet-sdk"

// 使用默认配置初始化，将从环境变量读取认证信息
client, err := intranet_sdk.NewClient()
if err != nil {
    // 处理错误
}
```

### 2.2 自定义配置初始化

```go
import "github.com/hujia-team/intranet-sdk"

// 使用配置选项自定义初始化
client, err := intranet_sdk.NewClient(
    // 自定义基础URL
    intranet_sdk.WithBaseURL("https://custom-api.example.com"),
    // 使用STS认证
    intranet_sdk.WithAccessKeyID("your_access_key_id"),
    intranet_sdk.WithAccessKeySecret("your_access_key_secret"),
    // 设置用户代理
    intranet_sdk.WithUserAgent("my-application/1.0"),
    // 启用调试模式
    intranet_sdk.WithDebug(true),
)
if err != nil {
    // 处理错误
}
```

## 3. 用户信息相关API

### 3.1 获取用户信息

```go
// 获取当前用户信息
userInfo, err := client.GetUserInfo()
if err != nil {
    // 处理错误
}

fmt.Printf("用户名: %s\n", userInfo.Name)
fmt.Printf("邮箱: %s\n", userInfo.Email)
fmt.Printf("角色: %s\n", userInfo.Role)
```

### 3.2 更新用户信息

```go
// 更新用户信息
updateRequest := &models.UserUpdateRequest{
    Name:  "新用户名",
    Email: "new.email@example.com",
}

updatedUser, err := client.UpdateUserInfo(updateRequest)
if err != nil {
    // 处理错误
}
```

## 4. 连接器相关API

### 4.1 获取连接器列表

```go
// 获取连接器列表
connectors, err := client.GetConnectors()
if err != nil {
    // 处理错误
}

// 遍历连接器列表
for _, connector := range connectors {
    fmt.Printf("连接器ID: %s\n", connector.ID)
    fmt.Printf("连接器名称: %s\n", connector.Name)
    fmt.Printf("连接器状态: %s\n", connector.Status)
}
```

### 4.2 获取连接器详情

```go
// 获取特定连接器详情
connectorID := "connector-123"
connector, err := client.GetConnectorByID(connectorID)
if err != nil {
    // 处理错误
}
```

### 4.3 创建新连接器

```go
// 创建新连接器
createRequest := &models.ConnectorCreateRequest{
    Name:        "新连接器",
    Description: "连接器描述信息",
    Type:        "database",
    Config: map[string]interface{}{
        "host":     "localhost",
        "port":     3306,
        "database": "test_db",
    },
}

newConnector, err := client.CreateConnector(createRequest)
if err != nil {
    // 处理错误
}
```

## 5. 错误处理

SDK使用自定义的错误类型，提供详细的错误信息。

### 5.1 基本错误处理

```go
// 基本错误处理
response, err := client.GetUserInfo()
if err != nil {
    fmt.Printf("操作失败: %v\n", err)
    return
}
```

### 5.2 高级错误处理（类型断言）

```go
// 高级错误处理，区分SDK错误和其他错误
response, err := client.GetUserInfo()
if err != nil {
    // 检查是否为SDK错误
    if sdkErr, ok := err.(*utils.SDKError); ok {
        // 访问SDK错误的详细信息
        fmt.Printf("错误码: %d\n", sdkErr.Code)
        fmt.Printf("错误信息: %s\n", sdkErr.Message)
        
        // 根据错误码进行不同的处理
        switch sdkErr.Code {
        case 401:
            fmt.Println("认证失败，请检查凭证")
        case 403:
            fmt.Println("权限不足")
        case 404:
            fmt.Println("资源不存在")
        default:
            fmt.Println("其他错误")
        }
    } else {
        // 处理非SDK错误（如网络错误等）
        fmt.Printf("系统错误: %v\n", err)
    }
    return
}
```

## 6. 日志和调试

### 6.1 启用调试模式

```go
// 启用调试模式
client, err := intranet_sdk.NewClient(
    intranet_sdk.WithDebug(true),
)
```

调试模式下，SDK会输出详细的请求和响应信息，包括：
- 请求URL
- 请求方法
- 请求头
- 请求体
- 响应状态码
- 响应头
- 响应体

### 6.2 自定义日志输出

```go
// 自定义日志函数
myLogger := func(level string, format string, args ...interface{}) {
    timestamp := time.Now().Format("2006-01-02 15:04:05")
    message := fmt.Sprintf(format, args...)
    fmt.Printf("[%s] [%s] %s\n", timestamp, level, message)
}

// 使用自定义日志函数
client, err := intranet_sdk.NewClient(
    intranet_sdk.WithLogger(myLogger),
)
```

## 7. 性能优化建议

### 7.1 重用客户端实例

在应用程序中，应该重用同一个客户端实例，而不是为每个请求创建新实例。这可以避免重复的初始化工作和资源消耗。

```go
// 应用程序启动时创建客户端实例（全局或单例）
var globalClient *intranet_sdk.Client

func initClient() error {
    var err error
    globalClient, err = intranet_sdk.NewClient()
    return err
}

// 在需要的地方使用全局客户端实例
func getUserInfo() (*models.UserInfo, error) {
    return globalClient.GetUserInfo()
}
```

### 7.2 设置超时

为API调用设置合理的超时时间，可以避免请求长时间挂起。

```go
// 设置超时时间
client, err := intranet_sdk.NewClient(
    intranet_sdk.WithTimeout(10 * time.Second),
)
```

## 8. 最佳实践

### 8.1 凭证管理

- 不要在代码中硬编码认证信息
- 使用环境变量或安全的配置管理系统
- 定期轮换凭证
- 使用最小权限原则

### 8.2 错误处理和重试

```go
// 带重试逻辑的API调用
func getWithRetry(client *intranet_sdk.Client, maxRetries int) (*models.UserInfo, error) {
    var userInfo *models.UserInfo
    var err error
    
    for i := 0; i < maxRetries; i++ {
        userInfo, err = client.GetUserInfo()
        if err == nil {
            return userInfo, nil
        }
        
        // 只有特定错误才重试
        if sdkErr, ok := err.(*utils.SDKError); ok {
            // 服务器错误可以重试
            if sdkErr.Code >= 500 && sdkErr.Code < 600 {
                fmt.Printf("重试第 %d 次...\n", i+1)
                time.Sleep(time.Duration(i+1) * 100 * time.Millisecond) // 指数退避
                continue
            }
        }
        
        // 其他错误不重试
        return nil, err
    }
    
    return nil, fmt.Errorf("达到最大重试次数: %w", err)
}
```

### 8.3 并发安全

客户端实例支持并发调用，但建议在并发环境中适当使用互斥锁来保护共享状态。

```go
// 并发安全的使用示例
func processBatch(client *intranet_sdk.Client, userIDs []string) []*models.UserInfo {
    var mu sync.Mutex
    var results []*models.UserInfo
    var wg sync.WaitGroup
    
    for _, userID := range userIDs {
        wg.Add(1)
        go func(id string) {
            defer wg.Done()
            
            userInfo, err := client.GetUserByID(id)
            if err != nil {
                fmt.Printf("获取用户 %s 信息失败: %v\n", id, err)
                return
            }
            
            mu.Lock()
            results = append(results, userInfo)
            mu.Unlock()
        }(userID)
    }
    
    wg.Wait()
    return results
}
```

## 9. 常见问题解答

### 9.1 认证失败

**问题**: 收到401错误，提示认证失败

**解决方案**:
- 检查AccessKeyID和AccessKeySecret是否正确
- 验证环境变量是否正确设置
- 确保STS令牌未过期

### 9.2 权限不足

**问题**: 收到403错误，提示权限不足

**解决方案**:
- 检查当前凭证是否有权限执行该操作
- 联系管理员申请相应权限

### 9.3 连接超时

**问题**: 请求超时

**解决方案**:
- 检查网络连接
- 验证API地址是否正确
- 尝试增加超时时间
- 检查服务器是否正常运行

## 10. 版本兼容性

SDK遵循语义化版本规范，版本号格式为MAJOR.MINOR.PATCH：

- 补丁版本（PATCH）：向后兼容的错误修复
- 小版本（MINOR）：向后兼容的功能新增
- 主版本（MAJOR）：不兼容的API更改

建议在生产环境中固定使用特定版本，并在升级前测试兼容性。

```bash
# 固定使用1.2.x系列的最新版本
go get github.com/hujia-team/intranet-sdk@v1.2.0
```

---

如需更多帮助，请查看示例代码或联系开发团队。