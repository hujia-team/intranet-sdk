# 架构说明

本文档描述当前 Go SDK 的代码组织和职责边界。

## 目录结构

```text
intranet.go              客户端入口
client/                  HTTP 客户端
models/                  数据模型
services/                各服务封装
utils/                   错误、日志、STS 等工具
examples/                示例
tests/                   集成测试
docs/                    文档
```

## 分层职责

### `intranet.go`

- 提供 `Client`
- 管理全局配置
- 暴露服务入口：
  - `User`
  - `Connector`
  - `ApiKey`
  - `Artifact`
  - `ClawSkill`

### `client/`

- 封装 HTTP 请求
- 统一处理认证头
- 统一处理错误状态码

认证优先级：

1. 登录态 Bearer Token
2. STS 认证
3. API Key

### `models/`

- 定义请求和响应结构
- 尽量贴近服务端 JSON 字段
- 承担少量通用辅助结构

### `services/`

- 面向业务场景封装接口
- 对原始接口做 helper 组合
- 典型例子：
  `GetChildArtifactHashesByCommitHash`
  它复用已有服务端接口，而不是要求服务端新增专用接口

### `tests/`

- `services/*_test.go`
  以 stub server 做单元测试
- `tests/*_integration_test.go`
  走真实环境联调

## 文档分层

- `README.md`
  入口页
- `usage_guide.md`
  通用使用指南
- `artifact_usage.md`
  制品专题
- `apikey-usage.md`
  ApiKey 专题
- `claw-skill-usage.md`
  本地 Skill 上传/重置 token 专题

## 设计原则

- 优先在 SDK 做 helper 组合，而不是为轻量消费场景新增服务端接口
- SDK 返回尽量保持 typed model
- 对需要保留状态码、原始响应文本和解析状态的能力，SDK 应返回 richer typed result，而不是把诊断信息留给调用方自己拼 HTTP
- 通用场景走基础接口，重复使用频率高的组合逻辑再沉到 helper
