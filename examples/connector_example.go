// 连接器示例
//
// 演示如何使用连接器相关功能，包括发送Kafka消息。
//
// 安装:
//
//	go get github.com/hujia-team/intranet-sdk
//
// 运行方法:
//
//	export INTRANET_ACCESS_KEY_ID=your_access_key_id
//	export INTRANET_ACCESS_KEY_SECRET=your_access_key_secret
//	go run examples/connector_example.go
package main

import (
	"fmt"
	"os"

	intranet "github.com/hujia-team/intranet-sdk"
)

type Message struct {
	Key string `json:"key"`
}

func main() {
	accessKeyID := os.Getenv("INTRANET_ACCESS_KEY_ID")
	accessKeySecret := os.Getenv("INTRANET_ACCESS_KEY_SECRET")

	if accessKeyID == "" || accessKeySecret == "" {
		fmt.Println("错误: 缺少STS凭证环境变量")
		fmt.Println("请设置 INTRANET_ACCESS_KEY_ID 和 INTRANET_ACCESS_KEY_SECRET 环境变量")
		os.Exit(1)
	}

	// 初始化SDK - 使用STS认证方式
	sdk, err := intranet.NewClient(
		intranet.WithAccessKeyID(accessKeyID),
		intranet.WithAccessKeySecret(accessKeySecret),
	)
	if err != nil {
		fmt.Printf("初始化SDK失败: %v\n", err)
		os.Exit(1)
	}
	resp, err := sdk.Connector.SendKafkaMessage("you-topic", Message{
		Key: "value",
	})
	if err != nil {
		fmt.Printf("发送消息到Kafka主题失败: %v\n", err)
	} else {
		fmt.Printf("错误码: %d\n", resp.Code)
		fmt.Printf("提示信息: %s\n", resp.Msg)
	}
}
