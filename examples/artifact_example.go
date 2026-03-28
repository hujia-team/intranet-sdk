package main

import (
	"fmt"
	"os"

	intranet "github.com/hujia-team/intranet-sdk"
	"github.com/hujia-team/intranet-sdk/models"
)

func main() {
	accessKeyID := os.Getenv("INTRANET_ACCESS_KEY_ID")
	accessKeySecret := os.Getenv("INTRANET_ACCESS_KEY_SECRET")

	if accessKeyID == "" || accessKeySecret == "" {
		fmt.Println("错误: 缺少STS凭证环境变量")
		fmt.Println("请设置 INTRANET_ACCESS_KEY_ID 和 INTRANET_ACCESS_KEY_SECRET 环境变量")
		os.Exit(1)
	}

	sdk, err := intranet.NewClient(
		intranet.WithAccessKeyID(accessKeyID),
		intranet.WithAccessKeySecret(accessKeySecret),
	)
	if err != nil {
		fmt.Printf("初始化SDK失败: %v\n", err)
		os.Exit(1)
	}

	list, err := sdk.Artifact.ListArtifacts(&models.ArtifactListReq{
		Page:     1,
		PageSize: 10,
	})
	if err != nil {
		fmt.Printf("获取制品列表失败: %v\n", err)
		return
	}
	fmt.Printf("制品总数: %d\n", list.Total)

	if len(list.Data) == 0 || list.Data[0].Name == nil {
		fmt.Println("没有可展示的制品")
		return
	}

	artifact, err := sdk.Artifact.GetArtifactByName(*list.Data[0].Name, nil)
	if err != nil {
		fmt.Printf("获取制品详情失败: %v\n", err)
		return
	}
	fmt.Printf("制品名称: %s\n", *artifact.Name)

	if artifact.CommitHash != nil {
		exists, err := sdk.Artifact.CheckExistsByCommitHash(*artifact.CommitHash, nil)
		if err != nil {
			fmt.Printf("按 commit hash 检查是否存在失败: %v\n", err)
			return
		}
		fmt.Printf("按 commit hash 检查存在性: %v\n", exists)

		plan, err := sdk.Artifact.PrepareDownloadByCommitHash(*artifact.CommitHash, nil, "./downloads")
		if err != nil {
			fmt.Printf("准备下载计划失败: %v\n", err)
			return
		}
		fmt.Printf("下载目标路径: %s\n", plan.TargetPath)

		metadata, err := sdk.Artifact.GetVersionMetadataByCommitHash(*artifact.CommitHash, nil)
		if err == nil && metadata != nil && metadata.MetadataFileName != nil {
			fmt.Printf("版本元数据文件: %s\n", *metadata.MetadataFileName)
		}
	}

	if artifact.ID != nil {
		download, err := sdk.Artifact.GetArtifactDownloadURL(*artifact.ID, "")
		if err != nil {
			fmt.Printf("获取下载地址失败: %v\n", err)
			return
		}
		fmt.Printf("下载文件名: %s\n", download.FileName)
		fmt.Printf("下载地址: %s\n", download.DownloadURL)
	}
}
