package main

import (
	"fmt"
	"log"
	"os"

	"github.com/hujia-team/intranet-sdk"
	"github.com/hujia-team/intranet-sdk/models"
	"github.com/hujia-team/intranet-sdk/utils"
)

func main() {
	// 从环境变量获取认证信息
	accessKeyID := os.Getenv("INTRANET_ACCESS_KEY_ID")
	accessKeySecret := os.Getenv("INTRANET_ACCESS_KEY_SECRET")

	if accessKeyID == "" || accessKeySecret == "" {
		log.Fatal("请设置环境变量: INTRANET_ACCESS_KEY_ID 和 INTRANET_ACCESS_KEY_SECRET")
	}

	// 启用调试模式
	utils.SetDefaultLogLevel(utils.LogLevelDebug)

	// 创建客户端实例
	client, err := intranet.NewClient(
		intranet.WithAccessKeyID(accessKeyID),
		intranet.WithAccessKeySecret(accessKeySecret),
	)
	if err != nil {
		log.Fatalf("创建客户端失败: %v", err)
	}

	fmt.Println("=== 演示 ApiKey 功能 ===\n")

	// 1. 获取 ApiKey 列表
	fmt.Println("1. 获取 ApiKey 列表...")
	listReq := &models.ApiKeyListReq{
		PageInfo: models.PageInfo{
			Page:     1,
			PageSize: 10,
		},
	}

	listResp, err := client.ApiKey.GetApiKeyList(listReq)
	if err != nil {
		log.Printf("获取 ApiKey 列表失败: %v", err)
	} else {
		fmt.Printf("   总数: %d\n", listResp.Total)
		for i, apiKey := range listResp.List {
			fmt.Printf("   [%d] ID: %d, Name: %s, Domain: %s\n", i+1, apiKey.ID, apiKey.Name, apiKey.Domain)
			if apiKey.BaseURL != "" {
				fmt.Printf("       BaseURL: %s\n", apiKey.BaseURL)
			}
			if apiKey.Description != "" {
				fmt.Printf("       Description: %s\n", apiKey.Description)
			}
			fmt.Printf("       权限 - Owner: %v, Admin: %v, Read: %v, Write: %v, Use: %v\n",
				apiKey.IsOwner, apiKey.IsAdmin, apiKey.HasRead, apiKey.HasWrite, apiKey.HasUse)
		}
	}
	fmt.Println()

	// 2. 如果列表中有数据,获取第一个 ApiKey 的详细信息
	if listResp != nil && len(listResp.List) > 0 {
		firstApiKeyID := listResp.List[0].ID
		fmt.Printf("2. 获取 ApiKey 详细信息 (ID: %d)...\n", firstApiKeyID)

		apiKeyInfo, err := client.ApiKey.GetApiKeyByID(firstApiKeyID)
		if err != nil {
			log.Printf("获取 ApiKey 详情失败: %v", err)
		} else {
			fmt.Printf("   ID: %d\n", apiKeyInfo.ID)
			fmt.Printf("   Name: %s\n", apiKeyInfo.Name)
			fmt.Printf("   Domain: %s\n", apiKeyInfo.Domain)
			if apiKeyInfo.BaseURL != "" {
				fmt.Printf("   BaseURL: %s\n", apiKeyInfo.BaseURL)
			}
			if apiKeyInfo.Description != "" {
				fmt.Printf("   Description: %s\n", apiKeyInfo.Description)
			}
			if apiKeyInfo.Sub != "" {
				fmt.Printf("   Sub: %s\n", apiKeyInfo.Sub)
			}
			if apiKeyInfo.SubType != "" {
				fmt.Printf("   SubType: %s\n", apiKeyInfo.SubType)
			}

			// 显示统计信息
			if apiKeyInfo.Stats != nil {
				fmt.Println("   统计信息:")
				fmt.Printf("     Active: %v\n", apiKeyInfo.Stats.IsActive)
				if apiKeyInfo.Stats.Usage != nil {
					fmt.Printf("     总使用 - Requests: %d, AllTokens: %d, Cost: %.4f\n",
						apiKeyInfo.Stats.Usage.Requests,
						apiKeyInfo.Stats.Usage.AllTokens,
						apiKeyInfo.Stats.Usage.Cost)
				}
			}
		}
		fmt.Println()
	}

	// 3. 创建新的 ApiKey (示例 - 需要根据实际情况调整)
	fmt.Println("3. 创建新的 ApiKey (示例)...")
	newApiKey := &models.ApiKeyInfo{
		Name:        "示例 ApiKey",
		Description: "这是一个通过 SDK 创建的示例 ApiKey",
		Domain:      "example-domain",
		BaseURL:     "https://api.example.com",
	}

	newID, err := client.ApiKey.CreateApiKey(newApiKey)
	if err != nil {
		log.Printf("创建 ApiKey 失败: %v (这是正常的,因为可能没有创建权限或域不存在)", err)
	} else {
		fmt.Printf("   创建成功! 新 ApiKey ID: %d\n", newID)

		// 4. 更新刚创建的 ApiKey
		fmt.Printf("4. 更新 ApiKey (ID: %d)...\n", newID)
		updateApiKey := &models.ApiKeyInfo{
			ID:          newID,
			Name:        "更新后的示例 ApiKey",
			Description: "描述已被更新",
		}

		err = client.ApiKey.UpdateApiKey(updateApiKey)
		if err != nil {
			log.Printf("更新 ApiKey 失败: %v", err)
		} else {
			fmt.Println("   更新成功!")
		}
		fmt.Println()

		// 5. 删除刚创建的 ApiKey
		fmt.Printf("5. 删除 ApiKey (ID: %d)...\n", newID)
		err = client.ApiKey.DeleteApiKey([]uint64{newID})
		if err != nil {
			log.Printf("删除 ApiKey 失败: %v", err)
		} else {
			fmt.Println("   删除成功!")
		}
	}

	fmt.Println("\n=== 演示完成 ===")
}
