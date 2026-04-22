package tests

import (
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/hujia-team/intranet-sdk/models"
	"github.com/hujia-team/intranet-sdk/utils"
)

func loadArtifactTestConfig(t *testing.T) (string, *models.ArtifactLookupOptions) {
	t.Helper()

	artifactName := os.Getenv("INTRANET_ARTIFACT_NAME")
	if artifactName == "" {
		t.Skip("未设置 INTRANET_ARTIFACT_NAME，跳过制品集成测试")
	}

	lookup := &models.ArtifactLookupOptions{}
	if modulePath := os.Getenv("INTRANET_ARTIFACT_MODULE_PATH"); modulePath != "" {
		lookup.ModulePath = modulePath
	}
	if artifactType := os.Getenv("INTRANET_ARTIFACT_TYPE"); artifactType != "" {
		lookup.ArtifactType = artifactType
	}
	if platform, ok := os.LookupEnv("INTRANET_ARTIFACT_PLATFORM"); ok {
		lookup.Platform = &platform
	}
	if semanticVersion := os.Getenv("INTRANET_ARTIFACT_SEMANTIC_VERSION"); semanticVersion != "" {
		lookup.SemanticVersion = semanticVersion
	}
	if includeVirtual := os.Getenv("INTRANET_ARTIFACT_INCLUDE_VIRTUAL"); includeVirtual != "" {
		value := includeVirtual == "1" || includeVirtual == "true" || includeVirtual == "TRUE"
		lookup.IncludeVirtual = &value
	}

	return artifactName, lookup
}

func TestArtifactReadFlow(t *testing.T) {
	utils.SetDefaultLogLevel(utils.LogLevelDebug)

	client, err := NewTestClient()
	if err != nil {
		t.Fatalf("创建客户端失败: %v", err)
	}

	artifactName, lookup := loadArtifactTestConfig(t)

	t.Log("=== 测试 Artifact Read Flow ===")
	t.Logf("目标制品名: %s", artifactName)

	listReq := &models.ArtifactListReq{
		Page:     1,
		PageSize: 20,
		Name:     &artifactName,
	}
	if lookup.ModulePath != "" {
		listReq.ModulePath = &lookup.ModulePath
	}
	if lookup.ArtifactType != "" {
		listReq.Type = &lookup.ArtifactType
	}
	if lookup.Platform != nil {
		listReq.Platform = lookup.Platform
	}
	if lookup.SemanticVersion != "" {
		listReq.SemanticVersion = &lookup.SemanticVersion
	}
	if lookup.IncludeVirtual != nil {
		listReq.IsVirtual = lookup.IncludeVirtual
	}

	t.Log("\n1. 查询制品列表...")
	listResp, err := client.Artifact.ListArtifacts(listReq)
	if err != nil {
		t.Fatalf("查询制品列表失败: %v", err)
	}
	if listResp == nil {
		t.Fatal("制品列表响应为空")
	}
	t.Logf("列表返回总数: %d, 当前页数量: %d", listResp.Total, len(listResp.Data))
	if len(listResp.Data) == 0 {
		t.Fatalf("未找到目标制品: %s", artifactName)
	}

	t.Log("\n2. 按名称获取制品详情...")
	artifact, err := client.Artifact.GetArtifactByName(artifactName, lookup)
	if err != nil {
		t.Fatalf("按名称获取制品详情失败: %v", err)
	}
	if artifact == nil {
		t.Fatal("制品详情为空")
	}
	if artifact.ID == nil {
		t.Fatal("制品 ID 为空")
	}
	if artifact.Name == nil || *artifact.Name == "" {
		t.Fatal("制品名称为空")
	}

	t.Logf("制品 ID: %d", *artifact.ID)
	t.Logf("制品名称: %s", *artifact.Name)
	if artifact.ProjectName != nil {
		t.Logf("项目名: %s", *artifact.ProjectName)
	}
	if artifact.ModulePath != nil {
		t.Logf("模块路径: %s", *artifact.ModulePath)
	}
	if artifact.SemanticVersion != nil {
		t.Logf("语义化版本: %s", *artifact.SemanticVersion)
	}

	t.Log("\n3. 按 ID 再次获取详情校验...")
	artifactByID, err := client.Artifact.GetArtifactByID(*artifact.ID)
	if err != nil {
		t.Fatalf("按 ID 获取制品详情失败: %v", err)
	}
	if artifactByID == nil || artifactByID.ID == nil || *artifactByID.ID != *artifact.ID {
		t.Fatalf("按 ID 获取到的制品不匹配: %#v", artifactByID)
	}

	if artifact.Tags != nil && *artifact.Tags != "" {
		t.Log("\n4. 获取并解析 tag schema...")
		version := ""
		if artifact.TagSchemaVersion != nil {
			version = *artifact.TagSchemaVersion
		}

		schema, err := client.Artifact.GetArtifactTagSchema(version)
		if err != nil {
			t.Fatalf("获取 tag schema 失败: %v", err)
		}
		if schema == nil {
			t.Fatal("tag schema 为空")
		}
		t.Logf("schema 版本: %s", schema.Version)

		parsedTags, err := client.Artifact.ParseArtifactTags(*artifact.Tags, schema)
		if err != nil {
			t.Fatalf("解析 tags 失败: %v", err)
		}
		t.Logf("解析后的 tags 字段数: %d", len(parsedTags))

		parsedByHelper, err := client.Artifact.GetParsedArtifactTags(*artifact.ID)
		if err != nil {
			t.Fatalf("通过 helper 获取解析后的 tags 失败: %v", err)
		}
		t.Logf("helper 返回的 tags 字段数: %d", len(parsedByHelper))
	} else {
		t.Log("\n4. 当前制品没有 tags，跳过 schema 与 tags 解析测试")
	}

	if artifact.ProjectName != nil && *artifact.ProjectName != "" {
		t.Log("\n5. 获取 JFrog token...")
		token, err := client.Artifact.GetJfrogToken(*artifact.ProjectName)
		if err != nil {
			t.Fatalf("获取 JFrog token 失败: %v", err)
		}
		if token == nil || token.AccessToken == "" {
			t.Fatalf("JFrog token 响应无效: %#v", token)
		}
		t.Logf("token 类型: %s, 过期秒数: %d", token.TokenType, token.ExpiresIn)

		tokenByName, err := client.Artifact.GetJfrogTokenByArtifactName(artifactName, lookup)
		if err != nil {
			t.Fatalf("按制品名获取 JFrog token 失败: %v", err)
		}
		if tokenByName == nil || tokenByName.AccessToken == "" {
			t.Fatalf("按制品名获取的 JFrog token 无效: %#v", tokenByName)
		}
	} else {
		t.Log("\n5. 当前制品没有 projectName，跳过 JFrog token 测试")
	}

	t.Log("\n6. 获取下载地址...")
	downloadURL, err := client.Artifact.GetArtifactDownloadURL(*artifact.ID, "")
	if err != nil {
		t.Fatalf("获取下载地址失败: %v", err)
	}
	if downloadURL == nil || downloadURL.DownloadURL == "" {
		t.Fatalf("下载地址响应无效: %#v", downloadURL)
	}
	t.Logf("下载文件名: %s", downloadURL.FileName)
	t.Logf("下载路径: %s", downloadURL.FilePath)

	downloadURLByName, err := client.Artifact.GetArtifactDownloadURLByName(artifactName, lookup, "")
	if err != nil {
		t.Fatalf("按制品名获取下载地址失败: %v", err)
	}
	if downloadURLByName == nil || downloadURLByName.DownloadURL == "" {
		t.Fatalf("按制品名获取下载地址响应无效: %#v", downloadURLByName)
	}

	if artifact.CommitHash != nil && *artifact.CommitHash != "" {
		t.Log("\n7. 检查高优先级 helper...")

		existsByCommitHash, err := client.Artifact.CheckExistsByCommitHash(*artifact.CommitHash, lookup)
		if err != nil {
			t.Fatalf("按 commit hash 检查制品是否存在失败: %v", err)
		}
		if !existsByCommitHash {
			t.Fatalf("按 commit hash 检查结果异常，期望存在: %s", *artifact.CommitHash)
		}
		t.Logf("CheckExistsByCommitHash: %v", existsByCommitHash)

		existsByName, err := client.Artifact.CheckExistsByName(artifactName, lookup)
		if err != nil {
			t.Fatalf("按名称检查制品是否存在失败: %v", err)
		}
		if !existsByName {
			t.Fatalf("按名称检查结果异常，期望存在: %s", artifactName)
		}
		t.Logf("CheckExistsByName: %v", existsByName)

		prepareDir := t.TempDir()
		plan, err := client.Artifact.PrepareDownloadByArtifactID(*artifact.ID, prepareDir)
		if err != nil {
			t.Fatalf("按 artifact ID 准备下载计划失败: %v", err)
		}
		if plan == nil || plan.Artifact == nil || plan.DownloadURL == nil || plan.Token == nil {
			t.Fatalf("下载计划响应无效: %#v", plan)
		}
		t.Logf("下载计划目标路径: %s", plan.TargetPath)
		t.Logf("下载计划文件路径: %s", plan.DownloadURL.FilePath)
		t.Logf("下载计划 checksum: %s", plan.Checksum)
		expectedTarget := filepath.Join(prepareDir, plan.DownloadURL.FileName)
		if plan.TargetPath != expectedTarget {
			t.Fatalf("下载计划目标路径不匹配: got=%s want=%s", plan.TargetPath, expectedTarget)
		}

		if os.Getenv("INTRANET_ARTIFACT_DOWNLOAD") == "true" {
			downloadedPlan, err := client.Artifact.DownloadByArtifactID(*artifact.ID, prepareDir)
			if err != nil {
				t.Fatalf("按 artifact ID 下载制品失败: %v", err)
			}
			if downloadedPlan.SkippedExisting {
				t.Fatal("首次下载不应命中 skip existing")
			}
			infoBefore, err := os.Stat(downloadedPlan.TargetPath)
			if err != nil {
				t.Fatalf("获取首次下载文件信息失败: %v", err)
			}
			t.Logf("首次下载完成: %s (%d bytes)", downloadedPlan.TargetPath, infoBefore.Size())

			time.Sleep(1100 * time.Millisecond)

			downloadedPlanAgain, err := client.Artifact.DownloadByArtifactID(*artifact.ID, prepareDir)
			if err != nil {
				t.Fatalf("重复按 artifact ID 下载制品失败: %v", err)
			}
			if !downloadedPlanAgain.SkippedExisting {
				t.Fatalf("重复下载未命中 skip existing: %#v", downloadedPlanAgain)
			}
			infoAfter, err := os.Stat(downloadedPlanAgain.TargetPath)
			if err != nil {
				t.Fatalf("获取重复下载文件信息失败: %v", err)
			}
			if !infoAfter.ModTime().Equal(infoBefore.ModTime()) {
				t.Fatalf("重复下载不应覆盖已有文件: before=%s after=%s", infoBefore.ModTime(), infoAfter.ModTime())
			}
			t.Logf("重复下载命中 skip existing: %s", downloadedPlanAgain.TargetPath)
		} else {
			t.Log("INTRANET_ARTIFACT_DOWNLOAD 未开启，跳过真实下载与 skip existing 验证")
		}

		if artifact.MetadataPath != nil && *artifact.MetadataPath != "" {
			metadata, err := client.Artifact.GetVersionMetadataByCommitHash(*artifact.CommitHash, lookup)
			if err != nil {
				t.Fatalf("按 commit hash 获取版本元数据失败: %v", err)
			}
			if metadata == nil || metadata.RawContent == nil || *metadata.RawContent == "" {
				t.Fatalf("版本元数据响应无效: %#v", metadata)
			}
			t.Logf("元数据文件: %s", valueOrFallback(metadata.MetadataFileName, "<unknown>"))
			t.Logf("元数据路径: %s", valueOrFallback(metadata.MetadataPath, "<unknown>"))
			t.Logf("元数据解析字段数: %d", len(metadata.Parsed))
		} else {
			t.Log("当前制品没有 metadataPath，跳过版本元数据 helper 测试")
		}
	} else {
		t.Log("\n7. 当前制品没有 commit hash，跳过高优先级 helper 测试")
	}

	t.Log("\n=== 测试完成 ===")
}

func TestArtifactChildHashesByCommitHash(t *testing.T) {
	utils.SetDefaultLogLevel(utils.LogLevelDebug)

	client, err := NewTestClient()
	if err != nil {
		t.Fatalf("创建客户端失败: %v", err)
	}

	testCases := []struct {
		name       string
		commitHash string
		lookup     *models.ArtifactLookupOptions
	}{
		{
			name:       "pkg root artifact",
			commitHash: "89a84fcee9c8db4c7d8ccb3547cfcc0a",
			lookup: &models.ArtifactLookupOptions{
				ArtifactType: "pkg",
			},
		},
		{
			name:       "app root artifact",
			commitHash: "60e604273314bafa694b30438bcc0c9a",
			lookup: &models.ArtifactLookupOptions{
				ArtifactType: "app",
			},
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Logf("commit hash: %s", tc.commitHash)

			rootArtifact, err := client.Artifact.GetArtifactByCommitHash(tc.commitHash, tc.lookup)
			if err != nil {
				t.Fatalf("通过 commit hash 获取根制品失败: %v", err)
			}
			if rootArtifact == nil || rootArtifact.ID == nil {
				t.Fatalf("根制品响应无效: %#v", rootArtifact)
			}

			childHashes, err := client.Artifact.GetChildArtifactHashesByCommitHash(tc.commitHash, tc.lookup)
			if err != nil {
				t.Fatalf("获取子制品 hashes 失败: %v", err)
			}
			if childHashes == nil {
				t.Fatal("子制品 hashes 响应为空")
			}

			t.Logf("root id: %d", *rootArtifact.ID)
			if rootArtifact.Name != nil {
				t.Logf("root name: %s", *rootArtifact.Name)
			}
			if rootArtifact.Type != nil {
				t.Logf("root type: %s", *rootArtifact.Type)
			}
			t.Logf("child hash count: %d", len(childHashes.ChildHashes))

			preview := make([]string, 0, min(5, len(childHashes.ChildHashes)))
			for i, item := range childHashes.ChildHashes {
				if i >= 5 {
					break
				}
				parts := make([]string, 0, 4)
				if item.Name != nil {
					parts = append(parts, "name="+*item.Name)
				}
				if item.Type != nil {
					parts = append(parts, "type="+*item.Type)
				}
				if item.CommitHash != nil {
					parts = append(parts, "hash="+*item.CommitHash)
				}
				if item.ParentID != nil {
					parts = append(parts, "parent="+uint64ToString(*item.ParentID))
				}
				preview = append(preview, strings.Join(parts, ", "))
			}
			if len(preview) > 0 {
				t.Logf("child hash preview:\n%s", strings.Join(preview, "\n"))
			}
		})
	}
}

func uint64ToString(v uint64) string {
	return strconv.FormatUint(v, 10)
}

func valueOrFallback(value *string, fallback string) string {
	if value == nil || *value == "" {
		return fallback
	}
	return *value
}
