#!/usr/bin/env python3
"""连接器示例

演示如何使用连接器相关功能，包括发送Kafka消息。

运行方法:
    # 方法1: 使用 .env 文件（推荐）
    # 1. 复制 .env.example 为 .env
    # 2. 编辑 .env 填入实际凭证
    # 3. 运行示例
    python examples/connector_example.py

    # 方法2: 使用环境变量
    export INTRANET_ACCESS_KEY_ID=your_access_key_id
    export INTRANET_ACCESS_KEY_SECRET=your_access_key_secret
    python examples/connector_example.py
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

    # 初始化SDK - 使用STS认证方式
    client_config = {
        "access_key_id": access_key_id,
        "access_key_secret": access_key_secret,
    }
    if base_url:
        client_config["base_url"] = base_url

    client = Client(**client_config)

    try:
        # 发送消息到Kafka主题
        print("\n正在发送消息到 Kafka 主题...")
        response = client.connector.send_kafka_message(
            topic="test",
            message={
                "key": "value",
                "key2": "value2",
                "timestamp": "2024-12-13",
            }
        )

        print("\n响应信息：")
        print("=" * 50)
        print(f"错误码   : {response.code}")
        print(f"提示信息 : {response.msg}")
        print("=" * 50)

        if response.is_success():
            print("✓ 消息发送成功!")
        else:
            print("❌ 消息发送失败!")
    except Exception as e:
        print(f"\n❌ 发送消息到Kafka主题失败: {e}")
        import traceback
        traceback.print_exc()
        sys.exit(1)


if __name__ == "__main__":
    main()
