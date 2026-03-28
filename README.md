# 内网客户端SDK (intranet-sdk)

内网客户端SDK是一个用于访问内网系统API的多语言客户端库，提供了简单、高效、可靠的方式来与内网服务进行交互。

## 安装

使用Go模块管理依赖，执行以下命令安装SDK：

```bash
go get github.com/hujia-team/intranet-sdk
```

## 快速开始

```go
import intranet "github.com/hujia-team/intranet-sdk"

sdk, err := intranet.NewClient(
	intranet.WithBaseURL("https://intranet.minieye.tech/sys-api"),
	intranet.WithAccessKeyID("your_access_key_id"),
	intranet.WithAccessKeySecret("your_access_key_secret"),
)
if err != nil {
	// 处理错误
}
```

## 文档索引

- 使用指南: [docs/usage_guide.md](docs/usage_guide.md)
- 制品能力: [docs/artifact_usage.md](docs/artifact_usage.md)
- ApiKey 使用: [docs/apikey-usage.md](docs/apikey-usage.md)
- 架构说明: [docs/architecture.md](docs/architecture.md)

## Artifact 高优先级能力

当前 Go SDK 已补齐一组面向 CLI/自动化调用的 artifact helper：

- 精确 `commit_hash` 定位制品
- 按 `commit_hash` / 名称检查制品是否存在
- 生成下载计划并复用 JFrog 原生下载能力
- 按 `commit_hash` 获取版本元数据并自动解析 JSON/XML

详细说明见 [docs/artifact_usage.md](docs/artifact_usage.md)。
