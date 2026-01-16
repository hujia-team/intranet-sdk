package tests

import (
	"testing"

	"github.com/hujia-team/intranet-sdk/models"
	"github.com/hujia-team/intranet-sdk/utils"
)

// TestGetCurrentGroup 测试获取当前订阅组功能
func TestGetCurrentGroup(t *testing.T) {
	// 启用调试模式
	utils.SetDefaultLogLevel(utils.LogLevelDebug)

	// 创建客户端实例
	client, err := NewTestClient()
	if err != nil {
		t.Fatalf("创建客户端失败: %v", err)
	}

	t.Log("=== 测试 GetCurrentGroup ===")

	// 调用 GetCurrentGroup
	currentGroup, err := client.ApiKey.GetCurrentGroup()
	if err != nil {
		t.Fatalf("获取当前订阅组失败: %v", err)
	}

	// 验证响应
	if currentGroup == nil {
		t.Fatal("响应为空")
	}

	// 打印结果
	if currentGroup.Group != nil {
		t.Logf("当前绑定的订阅组:")
		t.Logf("  ID: %d", currentGroup.Group.ID)
		t.Logf("  Name: %s", currentGroup.Group.Name)
		t.Logf("  Description: %s", currentGroup.Group.Description)

		// 验证必填字段
		if currentGroup.Group.ID == 0 {
			t.Error("分组 ID 不应为 0")
		}
		if currentGroup.Group.Name == "" {
			t.Error("分组名称不应为空")
		}
	} else {
		t.Log("当前未绑定任何订阅组")
	}

	t.Log("=== 测试完成 ===")
}

// TestGetAvailableGroups 测试获取可用订阅组列表
func TestGetAvailableGroups(t *testing.T) {
	// 启用调试模式
	utils.SetDefaultLogLevel(utils.LogLevelDebug)

	// 创建客户端实例
	client, err := NewTestClient()
	if err != nil {
		t.Fatalf("创建客户端失败: %v", err)
	}

	t.Log("=== 测试 GetAvailableGroups ===")

	// 调用 GetAvailableGroups
	availableGroups, err := client.ApiKey.GetAvailableGroups()
	if err != nil {
		t.Fatalf("获取可用订阅组失败: %v", err)
	}

	// 验证响应
	if availableGroups == nil {
		t.Fatal("响应为空")
	}

	t.Logf("可用订阅组数量: %d", len(availableGroups.Data))

	// 打印每个分组的详细信息
	for i, group := range availableGroups.Data {
		t.Logf("\n[%d] 订阅组信息:", i+1)
		t.Logf("  ID: %d", group.ID)
		t.Logf("  Name: %s", group.Name)
		t.Logf("  Description: %s", group.Description)

		// 配额信息

		t.Logf("  日配额:")
		t.Logf("    Limit: %d 分 (%.2f 美元)", group.DailyLimit, float64(group.DailyLimit)/100)
		t.Logf("    Used: %d 分 (%.2f 美元)", group.DailyUsed, float64(group.DailyUsed)/100)
		if group.DailyLimit > 0 {
			usagePercent := float64(group.DailyUsed) / float64(group.DailyLimit) * 100
			t.Logf("    Usage: %.2f%%", usagePercent)
		}
		usagePercent := float64(group.DailyUsed) / float64(group.DailyLimit) * 100
		t.Logf("    Usage: %.2f%%", usagePercent)

		t.Logf("  周配额:")
		t.Logf("    Limit: %d 分 (%.2f 美元)", group.WeeklyLimit, float64(group.WeeklyLimit)/100)
		t.Logf("    Used: %d 分 (%.2f 美元)", group.WeeklyUsed, float64(group.WeeklyUsed)/100)
		if group.WeeklyLimit > 0 {
			usagePercent := float64(group.WeeklyUsed) / float64(group.WeeklyLimit) * 100
			t.Logf("    Usage: %.2f%%", usagePercent)
		}

		// 验证必填字段
		if group.ID == 0 {
			t.Errorf("分组 [%d] ID 不应为 0", i+1)
		}
		if group.Name == "" {
			t.Errorf("分组 [%d] 名称不应为空", i+1)
		}
	}

	t.Log("\n=== 测试完成 ===")
}

// TestSwitchGroup 测试切换订阅组功能
func TestSwitchGroup(t *testing.T) {
	// 启用调试模式
	utils.SetDefaultLogLevel(utils.LogLevelDebug)

	// 创建客户端实例
	client, err := NewTestClient()
	if err != nil {
		t.Fatalf("创建客户端失败: %v", err)
	}

	t.Log("=== 测试 SwitchGroup ===")

	// 1. 先获取当前分组
	t.Log("\n1. 获取当前分组...")
	currentGroup, err := client.ApiKey.GetCurrentGroup()
	if err != nil {
		t.Fatalf("获取当前订阅组失败: %v", err)
	}

	var originalGroupID int64
	if currentGroup.Group != nil {
		originalGroupID = currentGroup.Group.ID
		t.Logf("   当前分组: %s (ID: %d)", currentGroup.Group.Name, currentGroup.Group.ID)
	} else {
		t.Log("   当前未绑定任何分组")
	}

	// 2. 获取可用分组列表
	t.Log("\n2. 获取可用分组列表...")
	availableGroups, err := client.ApiKey.GetAvailableGroups()
	if err != nil {
		t.Fatalf("获取可用订阅组失败: %v", err)
	}

	if len(availableGroups.Data) == 0 {
		t.Skip("没有可用的订阅组，跳过切换测试")
	}

	// 3. 选择一个不同的分组进行切换
	var targetGroupID int64
	var targetGroupName string
	for _, group := range availableGroups.Data {
		if group.ID != originalGroupID {
			targetGroupID = group.ID
			targetGroupName = group.Name
			break
		}
	}

	if targetGroupID == 0 {
		t.Skip("只有一个可用分组，无法测试切换功能")
	}

	t.Logf("   目标分组: %s (ID: %d)", targetGroupName, targetGroupID)

	// 4. 执行切换
	t.Log("\n3. 执行切换...")
	switchResp, err := client.ApiKey.SwitchGroup(&models.SwitchGroupReq{
		GroupID: targetGroupID,
	})
	if err != nil {
		t.Fatalf("切换订阅组失败: %v", err)
	}

	if !switchResp.Success {
		t.Fatalf("切换失败: %s", switchResp.Message)
	}

	t.Log("   切换成功!")
	if switchResp.Group != nil {
		t.Logf("   新分组: %s (ID: %d)", switchResp.Group.Name, switchResp.Group.ID)
	}

	// 5. 验证切换结果
	t.Log("\n4. 验证切换结果...")
	newCurrentGroup, err := client.ApiKey.GetCurrentGroup()
	if err != nil {
		t.Fatalf("获取当前订阅组失败: %v", err)
	}

	if newCurrentGroup.Group == nil {
		t.Fatal("切换后应该有绑定的分组")
	}

	if newCurrentGroup.Group.ID != targetGroupID {
		t.Errorf("切换后的分组 ID 不匹配: 期望 %d, 实际 %d", targetGroupID, newCurrentGroup.Group.ID)
	}

	t.Logf("   验证通过: 当前分组为 %s (ID: %d)", newCurrentGroup.Group.Name, newCurrentGroup.Group.ID)

	// 6. 切换回原来的分组（如果有）
	if originalGroupID != 0 {
		t.Log("\n5. 切换回原分组...")
		switchBackResp, err := client.ApiKey.SwitchGroup(&models.SwitchGroupReq{
			GroupID: originalGroupID,
		})
		if err != nil {
			t.Logf("   警告: 切换回原分组失败: %v", err)
		} else if switchBackResp.Success {
			t.Log("   已切换回原分组")
		}
	}

	t.Log("\n=== 测试完成 ===")
}
