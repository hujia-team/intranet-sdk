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
		fmt.Println("错误: 缺少 STS 凭证环境变量")
		return
	}

	sdk, err := intranet.NewClient(
		intranet.WithBaseURL("https://intranet.minieye.tech/sys-api"),
		intranet.WithAccessKeyID(accessKeyID),
		intranet.WithAccessKeySecret(accessKeySecret),
	)
	if err != nil {
		panic(err)
	}

	archive, err := os.ReadFile("./demo-skill.zip")
	if err != nil {
		panic(err)
	}

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

	fmt.Printf("status=%d\n", result.StatusCode)
	if result.Parsed != nil {
		fmt.Printf("msg=%s created=%v skipped=%v\n", result.Parsed.Msg, result.Parsed.Data.Created, result.Parsed.Data.Skipped)
		return
	}
	fmt.Printf("body=%s\n", result.BodyText)
}
