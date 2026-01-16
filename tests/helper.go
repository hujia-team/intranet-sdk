package tests

import (
	"fmt"
	"os"

	"github.com/hujia-team/intranet-sdk"
	"github.com/joho/godotenv"
)

// TestConfig 测试配置
type TestConfig struct {
	BaseURL         string
	AccessKeyID     string
	AccessKeySecret string
}

// LoadTestConfig 从 .env 文件加载测试配置
func LoadTestConfig() (*TestConfig, error) {
	// 尝试加载 .env 文件(如果存在)
	_ = godotenv.Load("../.env")

	config := &TestConfig{
		BaseURL:         os.Getenv("INTRANET_BASE_URL"),
		AccessKeyID:     os.Getenv("INTRANET_ACCESS_KEY_ID"),
		AccessKeySecret: os.Getenv("INTRANET_ACCESS_KEY_SECRET"),
	}

	// 验证必需的配置项
	if config.BaseURL == "" {
		return nil, fmt.Errorf("INTRANET_BASE_URL is required")
	}
	if config.AccessKeyID == "" {
		return nil, fmt.Errorf("INTRANET_ACCESS_KEY_ID is required")
	}
	if config.AccessKeySecret == "" {
		return nil, fmt.Errorf("INTRANET_ACCESS_KEY_SECRET is required")
	}

	return config, nil
}

// NewTestClient 创建用于测试的客户端
func NewTestClient() (*intranet.Client, error) {
	config, err := LoadTestConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to load test config: %w", err)
	}

	client, err := intranet.NewClient(
		intranet.WithBaseURL(config.BaseURL),
		intranet.WithAccessKeyID(config.AccessKeyID),
		intranet.WithAccessKeySecret(config.AccessKeySecret),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create client: %w", err)
	}

	return client, nil
}
