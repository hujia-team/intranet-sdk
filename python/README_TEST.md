# Python SDK 测试指南

## 快速开始测试

### 1. 配置环境变量

```bash
# 复制环境变量模板
cp .env.example .env

# 编辑 .env 文件，填入实际的凭证
vim .env  # 或使用其他编辑器
```

`.env` 文件内容示例：
```bash
INTRANET_ACCESS_KEY_ID=your_actual_access_key_id
INTRANET_ACCESS_KEY_SECRET=your_actual_access_key_secret
INTRANET_BASE_URL=https://intranet.minieye.tech/sys-api
INTRANET_SDK_LOG_LEVEL=DEBUG
```

### 2. 安装依赖

使用 uv（推荐）：
```bash
cd python
uv sync
```

或使用 pip：
```bash
cd python
pip install -e .
```

### 3. 运行测试

#### 方法 1: 运行完整测试脚本
```bash
cd python
python test_sdk.py
```

这会运行所有测试用例，包括：
- 获取用户信息
- 发送 Kafka 消息

#### 方法 2: 运行单个示例
```bash
# 测试获取用户信息
cd python
python examples/user_example.py

# 测试发送 Kafka 消息
python examples/connector_example.py
```

#### 方法 3: 使用 uv run（推荐）
```bash
cd python

# 运行测试脚本
uv run python test_sdk.py

# 或运行示例
uv run python examples/user_example.py
uv run python examples/connector_example.py
```

## 调试

### 启用 DEBUG 日志

在 `.env` 文件中设置：
```bash
INTRANET_SDK_LOG_LEVEL=DEBUG
```

或在代码中设置：
```python
from intranet_sdk import set_log_level

set_log_level("DEBUG")
```

### 时区问题

如果遇到时区错误，请确保系统时区设置为 UTC+8：

```bash
# Linux
sudo timedatectl set-timezone Asia/Shanghai

# macOS
sudo systemsetup -settimezone Asia/Shanghai

# 或在 Python 中临时设置
export TZ=Asia/Shanghai
```

## 常见问题

### 1. 导入错误
```
ModuleNotFoundError: No module named 'intranet_sdk'
```

**解决方法**：
```bash
# 安装 SDK
cd python
pip install -e .
```

### 2. 找不到 .env 文件
```
⚠ 未找到 .env 文件
```

**解决方法**：
```bash
cd python
cp .env.example .env
# 然后编辑 .env 填入实际凭证
```

### 3. STS 认证失败
```
❌ 错误: 系统时区必须为东八区(UTC+8)
```

**解决方法**：
设置系统时区为 Asia/Shanghai (UTC+8)

### 4. HTTP 错误
```
HTTP error: 401 Unauthorized
```

**可能原因**：
- Access Key ID 或 Secret 不正确
- STS token 生成失败（时区问题）
- 凭证已过期

**解决方法**：
- 检查 .env 文件中的凭证是否正确
- 检查系统时区
- 启用 DEBUG 日志查看详细信息

## 项目结构

```
python/
├── .env                    # 环境变量配置（需自行创建）
├── .env.example           # 环境变量模板
├── test_sdk.py            # 测试脚本
├── examples/              # 示例代码
│   ├── user_example.py
│   └── connector_example.py
├── src/intranet_sdk/      # SDK 源码
└── README_TEST.md         # 本文件
```

## 下一步

测试通过后，您可以：
1. 在自己的项目中使用 SDK
2. 发布到 PyPI：`uv publish`
3. 查看完整文档：[README.md](README.md)
