.PHONY: test test-unit test-integration help build release-patch release-minor release-major

# 获取当前版本
CURRENT_VERSION := $(shell git tag -l "v*" | sort -V | tail -n 1)
ifeq ($(CURRENT_VERSION),)
CURRENT_VERSION := v0.0.0
endif

# 默认目标
help:
	@echo "可用的 make 命令:"
	@echo ""
	@echo "测试相关:"
	@echo "  make test              - 运行所有测试（单元测试 + 集成测试）"
	@echo "  make test-unit         - 运行单元测试"
	@echo "  make test-integration  - 运行集成测试"
	@echo "  make test-verbose      - 运行所有测试（详细输出）"
	@echo ""
	@echo "构建相关:"
	@echo "  make build             - 编译检查所有代码"
	@echo "  make clean             - 清理构建文件"
	@echo ""
	@echo "发布相关:"
	@echo "  make release-patch     - 发布 patch 版本（bug 修复）"
	@echo "  make release-minor     - 发布 minor 版本（新功能）"
	@echo "  make release-major     - 发布 major 版本（破坏性更改）"
	@echo "  当前版本: $(CURRENT_VERSION)"
	@echo ""
	@echo "环境变量:"
	@echo "  GOWORK=off            - 禁用 Go workspace（保持 SDK 独立性）"

# 运行所有测试
test:
	@echo "运行所有测试..."
	@GOWORK=off go test -v $$(go list ./... | grep -v /examples)

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

# 编译检查
build:
	@echo "=== 编译检查 ==="
	@echo "检查代码..."
	@GOWORK=off go vet ./client ./models ./services ./utils .
	@echo "编译核心代码..."
	@GOWORK=off go build ./client ./models ./services ./utils .
	@echo "编译示例文件..."
	@cd examples && GOWORK=off go build user_example.go && GOWORK=off go build connector_example.go && cd ..
	@echo "✅ 编译检查通过"

# 清理构建文件
clean:
	@echo "清理构建文件..."
	@rm -f examples/user_example examples/connector_example
	@echo "✅ 清理完成"

# 版本升级函数
define bump_version
	$(eval VERSION_PARTS := $(subst ., ,$(subst v,,$(CURRENT_VERSION))))
	$(eval MAJOR := $(word 1,$(VERSION_PARTS)))
	$(eval MINOR := $(word 2,$(VERSION_PARTS)))
	$(eval PATCH := $(word 3,$(VERSION_PARTS)))
	$(if $(filter patch,$(1)),$(eval NEW_VERSION := v$(MAJOR).$(MINOR).$(shell echo $$(($(PATCH)+1)))))
	$(if $(filter minor,$(1)),$(eval NEW_VERSION := v$(MAJOR).$(shell echo $$(($(MINOR)+1))).0))
	$(if $(filter major,$(1)),$(eval NEW_VERSION := v$(shell echo $$(($(MAJOR)+1))).0.0))
endef

# 发布 patch 版本
release-patch:
	@echo "=== 发布 Patch 版本 ==="
	@echo "当前版本: $(CURRENT_VERSION)"
	$(call bump_version,patch)
	@echo "新版本: $(NEW_VERSION)"
	@$(MAKE) release-common VERSION=$(NEW_VERSION)

# 发布 minor 版本
release-minor:
	@echo "=== 发布 Minor 版本 ==="
	@echo "当前版本: $(CURRENT_VERSION)"
	$(call bump_version,minor)
	@echo "新版本: $(NEW_VERSION)"
	@$(MAKE) release-common VERSION=$(NEW_VERSION)

# 发布 major 版本
release-major:
	@echo "=== 发布 Major 版本 ==="
	@echo "当前版本: $(CURRENT_VERSION)"
	$(call bump_version,major)
	@echo "新版本: $(NEW_VERSION)"
	@$(MAKE) release-common VERSION=$(NEW_VERSION)

# 发布通用流程
release-common:
	@echo ""
	@echo "[步骤1] 检查 Git 状态..."
	@if ! git diff --quiet || ! git diff --cached --quiet; then \
		echo "❌ 错误: 存在未提交的更改"; \
		git status --short; \
		echo ""; \
		echo "请先提交所有更改后再发布版本"; \
		exit 1; \
	fi
	@echo "✅ Git 状态检查通过"
	@echo ""
	@echo "[步骤2] 运行测试..."
	@$(MAKE) test
	@echo ""
	@echo "[步骤3] 编译检查..."
	@$(MAKE) build
	@$(MAKE) clean
	@echo ""
	@echo "[步骤4] 整理依赖..."
	@GOWORK=off go mod tidy
	@GOWORK=off go mod verify
	@echo "✅ 依赖检查通过"
	@echo ""
	@echo "[步骤5] 创建 Git 标签..."
	@git tag -a "$(VERSION)" -m "Release $(VERSION)"
	@echo "✅ Git 标签 $(VERSION) 创建成功"
	@echo ""
	@echo "[步骤6] 推送到远程仓库..."
	@git push origin $$(git rev-parse --abbrev-ref HEAD)
	@git push origin "$(VERSION)"
	@echo "✅ 代码和标签已推送到远程仓库"
	@echo ""
	@echo "=== 发布完成 ==="
	@echo "版本: $(VERSION)"
	@echo "时间: $$(date)"
