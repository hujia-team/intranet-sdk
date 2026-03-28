# 制品能力使用说明

本文档只聚焦 `sdk.Artifact`。

## 能力列表

当前 Go SDK 已提供这些制品相关能力：

- `sdk.Artifact.ListArtifacts`
- `sdk.Artifact.GetArtifactByID`
- `sdk.Artifact.GetArtifactByName`
- `sdk.Artifact.GetArtifactByCommitHash`
- `sdk.Artifact.CheckExistsByCommitHash`
- `sdk.Artifact.CheckExistsByName`
- `sdk.Artifact.PrepareDownloadByCommitHash`
- `sdk.Artifact.DownloadByCommitHash`
- `sdk.Artifact.DownloadByName`
- `sdk.Artifact.GetVersionMetadataByCommitHash`
- `sdk.Artifact.GetChildArtifactHashesByCommitHash`
- `sdk.Artifact.GetArtifactTagSchema`
- `sdk.Artifact.ParseArtifactTags`
- `sdk.Artifact.GetParsedArtifactTags`
- `sdk.Artifact.GetJfrogToken`
- `sdk.Artifact.GetArtifactDownloadURL`
- `sdk.Artifact.GetArtifactDownloadURLByName`

## 按名称获取制品详情

```go
artifact, err := sdk.Artifact.GetArtifactByName(
	"D4Q2_V00.29.00_202603272037_release_snapshot_aarch64",
	&models.ArtifactLookupOptions{
		ArtifactType: "pkg",
	},
)
if err != nil {
	return err
}

fmt.Printf("artifact id: %d\n", *artifact.ID)
fmt.Printf("commit hash: %s\n", *artifact.CommitHash)
```

## 按 commit hash 获取根制品

```go
artifact, err := sdk.Artifact.GetArtifactByCommitHash(
	"89a84fcee9c8db4c7d8ccb3547cfcc0a",
	&models.ArtifactLookupOptions{
		ArtifactType: "pkg",
	},
)
if err != nil {
	return err
}

fmt.Printf("root artifact: %s (%s)\n", *artifact.Name, *artifact.Type)
```

注意：

- SDK 现在直接调用服务端精确 `commit_hash` 查询接口
- 建议传 `ArtifactType`
  常见值是 `pkg`、`app`

## 存在性检查

```go
exists, err := sdk.Artifact.CheckExistsByCommitHash(
	"89a84fcee9c8db4c7d8ccb3547cfcc0a",
	&models.ArtifactLookupOptions{
		ArtifactType: "pkg",
	},
)
if err != nil {
	return err
}

fmt.Printf("exists by commit hash: %v\n", exists)
```

也可以按名称做唯一性检查：

```go
exists, err := sdk.Artifact.CheckExistsByName(
	"D4Q2_V00.29.00_202603272037_release_snapshot_aarch64",
	&models.ArtifactLookupOptions{
		ArtifactType: "pkg",
	},
)
```

## 下载计划与下载

推荐先准备下载计划，再决定是否执行真正下载。

```go
plan, err := sdk.Artifact.PrepareDownloadByCommitHash(
	"89a84fcee9c8db4c7d8ccb3547cfcc0a",
	&models.ArtifactLookupOptions{
		ArtifactType: "pkg",
	},
	"./downloads",
)
if err != nil {
	return err
}

fmt.Printf("target path: %s\n", plan.TargetPath)
fmt.Printf("jfrog file path: %s\n", plan.DownloadURL.FilePath)
fmt.Printf("checksum: %s\n", plan.Checksum)
```

执行下载时，SDK 内部会复用 JFrog 原生 Go client：

```go
plan, err := sdk.Artifact.DownloadByCommitHash(
	"89a84fcee9c8db4c7d8ccb3547cfcc0a",
	&models.ArtifactLookupOptions{
		ArtifactType: "pkg",
	},
	"./downloads",
)
if err != nil {
	return err
}

fmt.Printf("downloaded to: %s\n", plan.TargetPath)
```

也可以按名称直接下载：

```go
_, err := sdk.Artifact.DownloadByName(
	"D4Q2_V00.29.00_202603272037_release_snapshot_aarch64",
	&models.ArtifactLookupOptions{
		ArtifactType: "pkg",
	},
	"./downloads",
)
```

## 版本元数据

```go
metadata, err := sdk.Artifact.GetVersionMetadataByCommitHash(
	"89a84fcee9c8db4c7d8ccb3547cfcc0a",
	&models.ArtifactLookupOptions{
		ArtifactType: "pkg",
	},
)
if err != nil {
	return err
}

fmt.Printf("metadata file: %s\n", *metadata.MetadataFileName)
fmt.Printf("raw size: %d\n", len(*metadata.RawContent))
fmt.Printf("parsed fields: %d\n", len(metadata.Parsed))
```

说明：

- `version.json` / `mcu_version.json` 会解析为 JSON map
- `bsp_version.xml` 会被转成通用 map 结构
- 原始内容始终保留在 `RawContent`

## 获取递归子制品 hashes

```go
childHashes, err := sdk.Artifact.GetChildArtifactHashesByCommitHash(
	"89a84fcee9c8db4c7d8ccb3547cfcc0a",
	&models.ArtifactLookupOptions{
		ArtifactType: "pkg",
	},
)
if err != nil {
	return err
}

fmt.Printf("child count: %d\n", len(childHashes.ChildHashes))
for _, item := range childHashes.ChildHashes {
	if item.Name != nil && item.CommitHash != nil {
		fmt.Printf("%s -> %s\n", *item.Name, *item.CommitHash)
	}
}
```

底层逻辑：

1. 调 `/aiplorer/artifact/by-commit-hash`
2. 服务端精确定位根制品
3. 再从详情里的递归 `dependencies` 提取所有子制品的 `commit_hash`

## 标签与 schema

```go
schema, err := sdk.Artifact.GetArtifactTagSchema("0.2.0")
if err != nil {
	return err
}

parsed, err := sdk.Artifact.ParseArtifactTags(rawTagsJSON, schema)
if err != nil {
	return err
}

fmt.Printf("parsed tag keys: %d\n", len(parsed))
```

也可以直接通过制品 ID 获取解析结果：

```go
parsed, err := sdk.Artifact.GetParsedArtifactTags(artifactID)
if err != nil {
	return err
}
```

## JFrog token 与下载地址

```go
token, err := sdk.Artifact.GetJfrogToken("D4Q2")
if err != nil {
	return err
}

downloadURL, err := sdk.Artifact.GetArtifactDownloadURL(artifactID, "artifact")
if err != nil {
	return err
}

fmt.Printf("jfrog token type: %s\n", token.TokenType)
fmt.Printf("download file: %s\n", downloadURL.FileName)
```

## 已验证样例

这些样例已经在正式服验证通过：

- `pkg` 根制品 `89a84fcee9c8db4c7d8ccb3547cfcc0a`
  根制品：`D4Q2_V00.29.00_202603272037_release_snapshot_aarch64`
  递归子制品 hash 数量：`20`
- `app` 根制品 `60e604273314bafa694b30438bcc0c9a`
  根制品：`app-apa/perception-1.0.0`
  递归子制品 hash 数量：`0`
