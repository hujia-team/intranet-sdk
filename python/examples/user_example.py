#!/usr/bin/env python3
"""用户信息示例

演示如何使用STS认证获取用户详细信息。

运行方法:
    # 方法1: 使用 .env 文件（推荐）
    # 1. 复制 .env.example 为 .env
    # 2. 编辑 .env 填入实际凭证
    # 3. 运行示例
    python examples/user_example.py

    # 方法2: 使用环境变量
    export INTRANET_ACCESS_KEY_ID=your_access_key_id
    export INTRANET_ACCESS_KEY_SECRET=your_access_key_secret
    python examples/user_example.py
"""

import os
import sys
from pathlib import Path

# 加载 .env 文件
try:
    from dotenv import load_dotenv
    # 尝试从 python 目录加载 .env
    env_path = Path(__file__).parent.parent / ".env"
    if env_path.exists():
        load_dotenv(env_path)
        print(f"✓ 已加载配置文件: {env_path}")
    else:
        print(f"ℹ 未找到 .env 文件，将使用环境变量")
except ImportError:
    print("ℹ python-dotenv 未安装，将使用环境变量")

from intranet_sdk import Client


def main():
    """主函数"""
    access_key_id = os.getenv("INTRANET_ACCESS_KEY_ID")
    access_key_secret = os.getenv("INTRANET_ACCESS_KEY_SECRET")
    base_url = os.getenv("INTRANET_BASE_URL")

    if not access_key_id or not access_key_secret:
        print("❌ 错误: 缺少STS凭证")
        print("\n请使用以下方式之一配置凭证：")
        print("1. 创建 .env 文件（推荐）：")
        print("   cp python/.env.example python/.env")
        print("   然后编辑 python/.env 填入实际凭证")
        print("\n2. 设置环境变量：")
        print("   export INTRANET_ACCESS_KEY_ID=your_access_key_id")
        print("   export INTRANET_ACCESS_KEY_SECRET=your_access_key_secret")
        sys.exit(1)

    # 创建客户端实例 - 使用STS认证
    client_config = {
        "access_key_id": access_key_id,
        "access_key_secret": access_key_secret,
    }
    if base_url:
        client_config["base_url"] = base_url

    client = Client(**client_config)

    try:
        # 获取用户信息
        print("\n正在获取用户信息...")
        user_info = client.user.get_user_info()

        print("\n✓ 成功获取用户信息：")
        print("=" * 50)
        print(f"用户名       : {user_info.username}")
        print(f"用户昵称     : {user_info.nickname}")
        print(f"用户ID       : {user_info.user_id}")
        print(f"用户头像     : {user_info.avatar}")
        print(f"主目录路径   : {user_info.home_path}")
        print(f"角色名称     : {user_info.role_name}")
        print(f"部门名称     : {user_info.department_name}")
        print("=" * 50)
    except Exception as e:
        print(f"\n❌ 获取用户信息失败: {e}")
        import traceback
        traceback.print_exc()
        sys.exit(1)


if __name__ == "__main__":
    main()
