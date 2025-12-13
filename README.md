# 内网客户端SDK (intranet-sdk)

内网客户端SDK是一个用于访问内网系统API的多语言客户端库，提供了简单、高效、可靠的方式来与内网服务进行交互。

## 支持的语言

- **Go SDK**: 位于项目根目录 ([快速开始](#go-sdk))
- **Python SDK**: 位于 `python/` 目录 ([Python 文档](python/README.md))

## 功能特性

- 简洁易用的API接口
- 支持STS认证方式
- 完整的错误处理机制
- 支持调试模式
- 可配置的基础参数

---

## Go SDK

### 安装

使用Go模块管理依赖，执行以下命令安装SDK：

```bash
go get github.com/hujia-team/intranet-sdk
```

### 快速开始

### 初始化客户端

```go
import "github.com/hujia-team/intranet-sdk"

// 创建客户端实例 - 使用STS认证
sdk, err := intranet.NewClient(
	intranet.WithBaseURL("https://intranet.minieye.tech/sys-api"),
	intranet.WithAccessKeyID("your_access_key_id"),
	intranet.WithAccessKeySecret("your_access_key_secret"),
)
if err != nil {
	// 处理错误
}
```

### 示例代码

Go SDK提供了示例代码，位于`examples`目录下。以下是各个示例文件的说明：

#### 1. 用户信息示例 (user_example.go)

演示如何使用STS认证获取用户详细信息。

**运行方法：**
```bash
# 设置环境变量
export INTRANET_ACCESS_KEY_ID=your_access_key_id
export INTRANET_ACCESS_KEY_SECRET=your_access_key_secret

# 运行示例
go run examples/user_example.go
```

#### 2. 连接器示例 (connector_example.go)

演示如何使用连接器相关功能。

**运行方法：**
```bash
# 设置环境变量
export INTRANET_ACCESS_KEY_ID=your_access_key_id
export INTRANET_ACCESS_KEY_SECRET=your_access_key_secret

# 运行示例
go run examples/connector_example.go
```

### 错误处理

SDK使用自定义的错误类型返回详细的错误信息：

```go
import "github.com/hujia-team/intranet-sdk/utils"

// 处理错误
if err != nil {
	// 直接处理错误
	fmt.Printf("错误: %v\n", err)
}

// 或者检查是否为SDK错误
if sdkErr, ok := err.(*utils.SDKError); ok {
	// 获取错误消息
	message := sdkErr.Error()
	// 处理SDK错误
}
```

### 配置选项

SDK使用Option模式进行配置：

| 配置方法 | 说明 |
|---------|------|
| WithBaseURL(url string) | 设置API基础URL，默认值: https://intranet.minieye.tech/sys-api |
| WithAuthToken(token string) | 设置访问令牌用于认证 |
| WithAPIKey(apiKey string) | 设置API密钥用于认证 |
| WithUserAgent(userAgent string) | 设置用户代理，默认值: minieye-intranet-sdk/1.0 |
| WithAccessKeyID(accessKeyID string) | 设置STS认证的AccessKeyID |
| WithAccessKeySecret(accessKeySecret string) | 设置STS认证的AccessKeySecret |

---

## Python SDK

Python SDK 使用 uv 包管理器构建，支持 Python 3.8+。

详细文档请查看: [Python SDK README](python/README.md)

### 快速安装

```bash
# 使用 pip
pip install intranet-sdk

# 使用 uv
uv add intranet-sdk
```

### 快速示例

```python
from intranet_sdk import Client

# 创建客户端
client = Client(
    access_key_id="your_access_key_id",
    access_key_secret="your_access_key_secret"
)

# 获取用户信息
user_info = client.user.get_user_info()
print(f"用户名: {user_info.username}")
```

---

## 注意事项

1. 妥善保管用户凭证，避免硬编码在代码中
2. 在生产环境中，建议通过环境变量或配置文件管理敏感信息
3. 实现适当的错误处理机制
4. STS凭证应定期轮换以保证安全性

## 发布与版本管理

本项目使用Git标签进行版本管理，遵循语义化版本规范。

### 发布流程

项目包含一个自动化发布脚本 `release.sh`，用于简化发布流程：

```bash
# 运行发布脚本
./release.sh v1.0.0  # 替换为实际版本号
```

发布脚本会执行以下操作：
1. 检查代码编译状态（核心代码和示例文件）
2. 清理构建文件
3. 运行 go mod tidy 检查依赖关系
4. 检查Git工作区状态
5. 创建并推送Git标签

### 版本号规范

版本号必须遵循语义化版本规范：
- MAJOR.MINOR.PATCH
- 修复错误但不改变API：增加PATCH版本
- 向后兼容的功能添加：增加MINOR版本
- 不向后兼容的API更改：增加MAJOR版本

## 许可证

[MIT License](LICENSE)

## 贡献

欢迎提交问题和拉取请求！

## 联系方式

如有任何问题，请联系开发团队：hujia-team@example.com