#!/bin/bash

# 内网客户端SDK发布脚本 - 支持自动版本升级
set -e

echo "=== 内网客户端SDK发布脚本 ==="

# 函数：获取当前最新版本号
get_latest_version() {
    # 尝试获取最新的git标签版本
    local latest_tag=$(git tag -l "v*" | sort -V | tail -n 1)
    if [ -z "$latest_tag" ]; then
        echo "v0.0.0"  # 如果没有标签，返回初始版本
    else
        echo "$latest_tag"
    fi
}

# 函数：自动升级版本号
bump_version() {
    local current_version=$1
    local bump_type=$2
    
    # 移除v前缀
    local version=${current_version:1}
    local major=$(echo $version | cut -d. -f1)
    local minor=$(echo $version | cut -d. -f2)
    local patch=$(echo $version | cut -d. -f3)
    
    case "$bump_type" in
        major)
            major=$((major + 1))
            minor=0
            patch=0
            ;;
        minor)
            minor=$((minor + 1))
            patch=0
            ;;
        patch)
            patch=$((patch + 1))
            ;;
        *)
            echo "错误: 不支持的升级类型。使用 major, minor 或 patch"
            exit 1
            ;;
    esac
    
    echo "v${major}.${minor}.${patch}"
}

# 获取当前最新版本
CURRENT_VERSION=$(get_latest_version)
echo "当前最新版本: $CURRENT_VERSION"

# 处理命令行参数
if [ -z "$1" ]; then
    # 如果没有提供参数，显示菜单让用户选择
    echo "
请选择要升级的版本部分:"
    echo "1) major - 不兼容的API更改"
    echo "2) minor - 向后兼容的功能添加"
    echo "3) patch - 向后兼容的错误修复"
    echo "4) 自定义版本号"
    
    read -p "请选择 (1-4): " choice
    
    case "$choice" in
        1)
            VERSION=$(bump_version "$CURRENT_VERSION" major)
            ;;
        2)
            VERSION=$(bump_version "$CURRENT_VERSION" minor)
            ;;
        3)
            VERSION=$(bump_version "$CURRENT_VERSION" patch)
            ;;
        4)
            read -p "请输入自定义版本号 (格式: v1.x.y): " custom_version
            if ! [[ "$custom_version" =~ ^v[0-9]+\.[0-9]+\.[0-9]+$ ]]; then
                echo "错误: 版本号格式不正确，请使用 v1.x.y 格式"
                exit 1
            fi
            VERSION=$custom_version
            ;;
        *)
            echo "错误: 无效的选择"
            exit 1
            ;;
    esac
else
    # 检查是否是bump命令
    if [[ "$1" == "bump" && "$2" =~ ^(major|minor|patch)$ ]]; then
        VERSION=$(bump_version "$CURRENT_VERSION" "$2")
    else
        # 直接使用提供的版本号
        VERSION=$1
    fi
fi

# 检查是否是--dry-run模式
DRY_RUN=false
if [[ "$VERSION" == "--dry-run" ]]; then
    DRY_RUN=true
    # 如果是--dry-run模式，取下一个参数作为版本或bump命令
    if [ -n "$2" ]; then
        if [[ "$2" == "bump" && "$3" =~ ^(major|minor|patch)$ ]]; then
            VERSION=$(bump_version "$CURRENT_VERSION" "$3")
        else
            VERSION=$2
        fi
    else
        echo "错误: --dry-run模式需要提供版本号或bump命令"
        exit 1
    fi
fi

# 检查版本号格式
if ! [[ "$VERSION" =~ ^v[0-9]+\.[0-9]+\.[0-9]+$ ]]; then
    echo "错误: 版本号格式不正确，请使用 v1.x.y 格式"
    exit 1
fi

echo "发布版本: $VERSION"
if [ "$DRY_RUN" = true ]; then
    echo "[DRY-RUN模式] 将不会创建实际的Git标签或推送更改"
fi

# 1. 检查代码编译状态 - 只编译SDK核心代码，排除examples
echo "
[步骤1] 检查代码编译状态..."
# 只检查和编译SDK核心代码
go vet ./client ./models ./services ./utils .
go build ./client ./models ./services ./utils .
# 分别编译每个示例文件，避免同时编译多个main函数
echo "编译示例文件..."
cd examples && go build user_example.go && go build connector_example.go && cd ..

# 2. 清理构建文件
echo "
[步骤2] 清理构建文件..."
rm -f examples/user_example examples/connector_example

# 3. 确保依赖关系正确
echo "
[步骤3] 整理依赖关系..."
go mod tidy
go mod verify

# 4. 检查Git状态
echo "
[步骤4] 检查Git状态..."
if ! git diff --quiet; then
    echo "警告: 存在未提交的更改，请先提交或隐藏更改"
    read -p "是否继续发布？(y/N): " confirm
    if [[ ! "$confirm" =~ ^[Yy]$ ]]; then
        echo "发布已取消"
        exit 1
    fi
fi

# 5. 创建Git标签（可选）
echo "
[步骤5] 创建Git标签..."
if [ "$DRY_RUN" = true ]; then
    echo "[DRY-RUN] 将要执行: git tag -a \"$VERSION\" -m \"Release $VERSION\""
else
    git tag -a "$VERSION" -m "Release $VERSION"
    echo "Git标签 $VERSION 创建成功"
fi

# 6. 推送标签到远程仓库（可选）
echo "
[步骤6] 推送标签到远程仓库..."
if [ "$DRY_RUN" = true ]; then
    echo "[DRY-RUN] 跳过交互式推送确认"
    echo "[DRY-RUN] 模拟将要执行: git push origin \"$VERSION\""
else
    read -p "是否推送标签到远程仓库？(y/N): " push_tag
    if [[ "$push_tag" =~ ^[Yy]$ ]]; then
        git push origin "$VERSION"
        echo "标签 $VERSION 已推送至远程仓库"
    fi
fi

# 7. 验证模块发布
echo "
[步骤7] 验证模块可用性..."
go list -m github.com/hujia-team/intranet-sdk@latest || echo "注意: 模块尚未公开发布"

echo "
=== 发布准备完成 ==="
echo "要公开发布此模块，请运行: go list -m github.com/hujia-team/intranet-sdk@$VERSION"
echo "请确保您有适当的权限发布到GitHub或Go模块代理"
echo "发布版本: $VERSION"
echo "发布时间: $(date)"