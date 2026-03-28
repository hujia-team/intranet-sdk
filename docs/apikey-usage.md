# ApiKey 能力使用说明

本文档只聚焦 `sdk.ApiKey`。

## 能力列表

当前 Go SDK 已提供这些 ApiKey 相关能力：

- `sdk.ApiKey.CreateApiKey`
- `sdk.ApiKey.UpdateApiKey`
- `sdk.ApiKey.DeleteApiKey`
- `sdk.ApiKey.GetApiKeyList`
- `sdk.ApiKey.GetApiKeyByID`
- `sdk.ApiKey.GetSub2ApiKey`
- `sdk.ApiKey.GetAvailableGroups`
- `sdk.ApiKey.GetCurrentGroup`
- `sdk.ApiKey.SwitchGroup`

## 获取 ApiKey 列表

```go
req := &models.ApiKeyListReq{
	PageInfo: models.PageInfo{
		Page:     1,
		PageSize: 10,
	},
}

resp, err := sdk.ApiKey.GetApiKeyList(req)
if err != nil {
	return err
}

fmt.Printf("total: %d\n", resp.Total)
```

## 获取单个 ApiKey

```go
apiKey, err := sdk.ApiKey.GetApiKeyByID(123)
if err != nil {
	return err
}

fmt.Printf("name: %s\n", apiKey.Name)
fmt.Printf("domain: %s\n", apiKey.Domain)
```

## 创建 / 更新 / 删除

```go
id, err := sdk.ApiKey.CreateApiKey(&models.ApiKeyInfo{
	Name:        "demo",
	Description: "for testing",
})
if err != nil {
	return err
}

err = sdk.ApiKey.UpdateApiKey(&models.ApiKeyInfo{
	ID:   id,
	Name: "demo-updated",
})
if err != nil {
	return err
}

err = sdk.ApiKey.DeleteApiKey([]uint64{id})
if err != nil {
	return err
}
```

## Sub2Api 分组

```go
currentGroup, err := sdk.ApiKey.GetCurrentGroup()
if err != nil {
	return err
}
_ = currentGroup

availableGroups, err := sdk.ApiKey.GetAvailableGroups()
if err != nil {
	return err
}
_ = availableGroups
```

说明：

- `GetCurrentGroup` / `SwitchGroup` 依赖账号具备对应的 `sub2api` 权限
- 无权限时服务端可能返回 `401` 或业务错误

## 参考

- 集成测试: [tests/apikey_integration_test.go](../tests/apikey_integration_test.go)
