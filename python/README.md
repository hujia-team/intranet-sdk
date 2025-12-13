# Intranet SDK - Python

MINIEYE 内网系统 API 的 Python 客户端库，提供了简单、高效、可靠的方式来与内网服务进行交互。

## 功能特性

- 简洁易用的 API 接口
- 支持 STS 认证方式
- 完整的错误处理机制
- 支持调试模式
- 可配置的基础参数
- 类型提示支持（Python 3.8+）

## 安装

### 使用 pip 安装

```bash
pip install intranet-sdk --index-url https://pypi.minieye.tech/
```

### 使用 uv 安装

```bash
uv add intranet-sdk --index-url https://pypi.minieye.tech/
```

### 从源码安装

```bash
# 克隆仓库
git clone https://github.com/hujia-team/intranet-sdk.git
cd intranet-sdk/python

# 使用 uv 安装依赖
uv sync

# 或使用 pip 安装
pip install -e .
```

## 快速开始

### 初始化客户端

```python
from intranet_sdk import Client

# 创建客户端实例 - 使用 STS 认证
client = Client(
    access_key_id="your_access_key_id",
    access_key_secret="your_access_key_secret"
)
```

### 获取用户信息

```python
# 获取当前用户信息
user_info = client.user.get_user_info()

print(f"用户名: {user_info.username}")
print(f"昵称: {user_info.nickname}")
print(f"部门: {user_info.department_name}")
```

### 发送 Kafka 消息

```python
# 发送消息到 Kafka 主题
response = client.connector.send_kafka_message(
    topic="test",
    message={
        "key": "value",
        "data": "hello world"
    }
)

if response.is_success():
    print("消息发送成功!")
else:
    print(f"消息发送失败: {response.msg}")
```

## 示例代码

SDK 提供了示例代码，位于 `examples` 目录下：

### 1. 用户信息示例 (`user_example.py`)

演示如何使用 STS 认证获取用户详细信息。

**运行方法：**
```bash
# 设置环境变量
export INTRANET_ACCESS_KEY_ID=your_access_key_id
export INTRANET_ACCESS_KEY_SECRET=your_access_key_secret

# 运行示例
python examples/user_example.py
```

### 2. 连接器示例 (`connector_example.py`)

演示如何使用连接器相关功能。

**运行方法：**
```bash
# 设置环境变量
export INTRANET_ACCESS_KEY_ID=your_access_key_id
export INTRANET_ACCESS_KEY_SECRET=your_access_key_secret

# 运行示例
python examples/connector_example.py
```

## 配置选项

Client 类支持以下配置参数：

| 参数 | 类型 | 说明 | 默认值 |
|------|------|------|--------|
| `base_url` | str | API 基础 URL | `https://intranet.minieye.tech/sys-api` |
| `api_key` | str | API 密钥（可选） | `None` |
| `access_key_id` | str | STS 认证的 Access Key ID（可选） | `None` |
| `access_key_secret` | str | STS 认证的 Access Key Secret（可选） | `None` |
| `user_agent` | str | 用户代理字符串 | `minieye-intranet-sdk-python/0.1.0` |
| `timeout` | int | 请求超时时间（秒） | `30` |

## 错误处理

SDK 使用自定义的错误类型返回详细的错误信息：

```python
from intranet_sdk import Client, APIError, InternalError

client = Client(
    access_key_id="your_key_id",
    access_key_secret="your_key_secret"
)

try:
    user_info = client.user.get_user_info()
    print(f"用户名: {user_info.username}")
except APIError as e:
    # API 返回的错误
    print(f"API 错误: {e}")
except InternalError as e:
    # SDK 内部错误（如网络错误、JSON 解析错误等）
    print(f"内部错误: {e}")
```

### 错误类型

- `SDKError`: SDK 所有错误的基类
- `APIError`: API 返回错误时抛出
- `InternalError`: SDK 内部错误（网络错误、序列化错误等）

## 日志配置

SDK 内置日志支持，可以通过环境变量或代码配置日志级别：

### 通过环境变量

```bash
export INTRANET_SDK_LOG_LEVEL=DEBUG
python your_script.py
```

### 通过代码

```python
from intranet_sdk import set_log_level

# 设置日志级别
set_log_level("DEBUG")  # 可选: DEBUG, INFO, WARNING, ERROR, CRITICAL
```

## API 参考

### Client

主客户端类，提供对所有服务的访问。

#### 属性

- `user`: 用户服务实例
- `connector`: 连接器服务实例

### UserService

用户相关操作服务。

#### 方法

- `get_user_info() -> UserInfo`: 获取当前用户信息

### ConnectorService

连接器相关操作服务。

#### 方法

- `send_kafka_message(topic: str, message: Any) -> BaseMsgResp`: 发送消息到 Kafka 主题

### 数据模型

#### UserInfo

用户信息模型。

**属性：**
- `user_id`: 用户唯一标识
- `username`: 用户名
- `nickname`: 昵称
- `avatar`: 头像 URL
- `home_path`: 主目录路径
- `role_name`: 角色名称
- `department_name`: 部门名称

#### BaseMsgResp

基础响应模型。

**属性：**
- `code`: 错误码（0 表示成功）
- `msg`: 响应消息

**方法：**
- `is_success() -> bool`: 判断响应是否成功

## 注意事项

1. 妥善保管用户凭证，避免硬编码在代码中
2. 在生产环境中，建议通过环境变量或配置文件管理敏感信息
3. 实现适当的错误处理机制
4. STS 凭证应定期轮换以保证安全性

## 开发

### 环境设置

```bash
# 克隆仓库
git clone https://github.com/hujia-team/intranet-sdk.git
cd intranet-sdk/python

# 使用 uv 安装开发依赖
uv sync --dev

# 或使用 pip
pip install -e ".[dev]"
```

### 运行测试

```bash
# 使用 pytest 运行测试
pytest

# 运行测试并显示覆盖率
pytest --cov=intranet_sdk
```

### 代码格式化

```bash
# 使用 ruff 格式化代码
ruff format .

# 使用 ruff 检查代码
ruff check .
```

## 许可证

MIT License

## 贡献

欢迎提交问题和拉取请求！

## 联系方式

如有任何问题，请联系开发团队或在 GitHub 上提交 Issue。

## 更新日志

### v0.1.0 (2024-12-13)

- 初始版本发布
- 支持 STS 认证
- 支持用户信息查询
- 支持 Kafka 消息发送
