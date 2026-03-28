# 使用指南

本文档只介绍 `intranet-sdk` 的通用使用方式，不展开某个具体服务的细节。

## 安装

```bash
go get github.com/hujia-team/intranet-sdk
```

## 初始化

推荐使用 STS 凭证初始化：

```go
import intranet "github.com/hujia-team/intranet-sdk"

sdk, err := intranet.NewClient(
	intranet.WithBaseURL("https://intranet.minieye.tech/sys-api"),
	intranet.WithAccessKeyID("your_access_key_id"),
	intranet.WithAccessKeySecret("your_access_key_secret"),
)
if err != nil {
	panic(err)
}
```

可用配置项：

- `WithBaseURL`
- `WithAPIKey`
- `WithUserAgent`
- `WithAccessKeyID`
- `WithAccessKeySecret`

## 服务入口

客户端当前暴露这些服务：

- `sdk.User`
- `sdk.Connector`
- `sdk.ApiKey`
- `sdk.Artifact`

## 最小示例

```go
userInfo, err := sdk.User.GetUserInfo()
if err != nil {
	return err
}

if userInfo.Username != nil {
	fmt.Printf("username: %s\n", *userInfo.Username)
}
```

## 错误处理

SDK 返回的主要是 `*utils.SDKError`：

```go
resp, err := sdk.User.GetUserInfo()
if err != nil {
	if sdkErr, ok := err.(*utils.SDKError); ok {
		fmt.Printf("sdk error: %s\n", sdkErr.Error())
	}
	return err
}
_ = resp
```

## 按场景查看专题文档

- 制品能力: [artifact_usage.md](./artifact_usage.md)
- ApiKey 能力: [apikey-usage.md](./apikey-usage.md)
- 架构说明: [architecture.md](./architecture.md)
