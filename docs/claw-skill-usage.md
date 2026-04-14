# Claw Skill 使用指南

本文档介绍 `intranet-sdk` 中本地 Skill 上传与重置上传 token 的能力。

## 适用接口

当前 SDK 暴露：

- `sdk.ClawSkill.UploadLocalSkill`
- `sdk.ClawSkill.ResetLocalSkillUploadToken`
- `sdk.ClawSkill.ReportPrivateSkillHubEvent`

这两个接口对应服务端：

- `POST /claw/skill/local/upload`
- `POST /claw/skill/local/token/reset`
- `POST /claw/skill/private-hub/event/report`

## 初始化

推荐使用 STS 凭证：

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

## 上传本地 Skill

```go
archive := []byte("zip-bytes")

result, err := sdk.ClawSkill.UploadLocalSkill(
	"https://intranet.minieye.tech/sys-api/claw/skill/local/upload",
	"demo-skill.zip",
	archive,
	"",
	"",
	nil,
)
if err != nil {
	panic(err)
}

if result.Parsed != nil {
	println(result.Parsed.Msg)
	println(result.Parsed.Data.Created)
	println(result.Parsed.Data.Skipped)
}
```

参数说明：

- 第 1 个参数支持完整上传 URL
- 第 2 个参数是上传文件名
- 第 3 个参数是 zip 包字节流
- 第 4 个参数是可选版本覆盖
- 第 5 个参数是更新已有 Skill 时的 `uploadToken`
- 第 6 个参数是额外请求头

返回结果说明：

- `StatusCode`: HTTP 状态码
- `BodyText`: 原始响应文本
- `Parsed`: 成功解析后的结构化响应
- `ParseError`: 响应不是预期 JSON 时的解析错误

## 重置上传 Token

```go
result, err := sdk.ClawSkill.ResetLocalSkillUploadToken(
	"https://intranet.minieye.tech/sys-api/claw/skill/local/token/reset",
	"demo-skill",
	nil,
)
if err != nil {
	panic(err)
}

if result.Parsed != nil {
	println(result.Parsed.Data.UploadToken)
}
```

## 上报 Private Skill Hub 事件

```go
skillName := "common.skill-hub.publisher"
action := "install"
success := true
clientName := "ai-forge"

result, err := sdk.ClawSkill.ReportPrivateSkillHubEvent(
	"https://intranet.minieye.tech/sys-api/claw/skill/private-hub/event/report",
	&models.PrivateSkillHubEventReportRequest{
		SkillName:  skillName,
		Action:     action,
		Success:    success,
		ClientName: &clientName,
	},
	nil,
)
if err != nil {
	panic(err)
}

if result.Parsed != nil {
	println(result.Parsed.Msg)
}
```

适用场景：

- `ai-forge` 下载或安装 `skill-hub` skill 后做 best-effort 上报
- 内部工具记录 skill 使用统计

参数说明：

- `SkillName`：skill 名称
- `Action`：`download` 或 `install`
- `Success`：是否成功
- 其余字段均为可选补充信息，如版本、操作人、客户端版本、错误摘要

## 服务端去重语义

本地 Skill 上传去重由服务端决定：

- 测试服/正式服建议关闭 `ForceLocalSkillUpload`
  这时服务端会计算并持久化内容 hash，无变更时返回 `skipped=true`
- 本地 9100 建议开启 `ForceLocalSkillUpload`
  这时服务端始终覆盖上传，不写数据库 hash

这意味着 SDK 调用方不需要自行实现服务端去重逻辑，只需消费返回结果中的：

- `created`
- `skipped`

## 建议

- CLI 或发布工具应优先读取 `StatusCode`、`BodyText`、`Parsed`
- 若 `Parsed == nil`，应直接输出 `BodyText` 便于排障
- 若更新已有 Skill，需妥善保存并传入 `uploadToken`
