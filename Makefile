.PHONY: test test-unit test-integration help

# 默认目标
help:
	@echo "可用的 make 命令:"
	@echo "  make test              - 运行所有测试（单元测试 + 集成测试）"
	@echo "  make test-unit         - 运行单元测试"
	@echo "  make test-integration  - 运行集成测试"
	@echo "  make test-verbose      - 运行所有测试（详细输出）"
	@echo ""
	@echo "环境变量:"
	@echo "  GOWORK=off            - 禁用 Go workspace（保持 SDK 独立性）"

# 运行所有测试
test:
	@echo "运行所有测试..."
	@GOWORK=off go test -v ./...

# 运行单元测试（排除 tests 目录）
test-unit:
	@echo "运行单元测试..."
	@GOWORK=off go test -v $$(go list ./... | grep -v /tests)

# 运行集成测试
test-integration:
	@echo "运行集成测试..."
	@if [ ! -f .env ]; then \
		echo "错误: .env 文件不存在"; \
		echo "请复制 .env.example 为 .env 并填入配置"; \
		exit 1; \
	fi
	@GOWORK=off go test -v ./tests -run Integration

# 运行所有测试（详细输出）
test-verbose:
	@echo "运行所有测试（详细输出）..."
	@GOWORK=off go test -v -count=1 ./...

# 运行特定的集成测试
test-current-group:
	@echo "运行 GetCurrentGroup 集成测试..."
	@GOWORK=off go test -v ./tests -run TestGetCurrentGroup

test-available-groups:
	@echo "运行 GetAvailableGroups 集成测试..."
	@GOWORK=off go test -v ./tests -run TestGetAvailableGroups

test-switch-group:
	@echo "运行 SwitchGroup 集成测试..."
	@GOWORK=off go test -v ./tests -run TestSwitchGroup
