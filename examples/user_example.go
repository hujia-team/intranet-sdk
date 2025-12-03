package main

import (
	"fmt"
	"os"

	intranet "github.com/hujia-team/intranet-sdk"
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
	resp, err := sdk.User.GetUserInfo()
	if err != nil {
		fmt.Printf("获取用户信息失败: %v\n", err)
	} else {
		fmt.Printf("用户信息: %v\n", resp.Username)
		fmt.Printf("用户昵称: %v\n", resp.Nickname)
		fmt.Printf("用户头像: %v\n", resp.Avatar)
		fmt.Printf("用户主目录路径: %v\n", resp.HomePath)
		fmt.Printf("用户角色名称: %v\n", resp.RoleName)
		fmt.Printf("用户部门名称: %v\n", resp.DepartmentName)
	}
}
