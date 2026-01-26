#!/usr/bin/env python3
"""连接器示例

演示如何使用连接器相关功能，包括发送Kafka消息。

安装:
    pip install intranet-sdk

运行方法:
    export INTRANET_ACCESS_KEY_ID=your_access_key_id
    export INTRANET_ACCESS_KEY_SECRET=your_access_key_secret
    python examples/connector_example.py
"""

import os
import sys

from intranet_sdk import Client


def main():
    """主函数"""
    access_key_id = os.getenv("INTRANET_ACCESS_KEY_ID")
    access_key_secret = os.getenv("INTRANET_ACCESS_KEY_SECRET")
    base_url = os.getenv("INTRANET_BASE_URL")

    if not access_key_id or not access_key_secret:
        print("错误: 缺少STS凭证环境变量")
        print("请设置 INTRANET_ACCESS_KEY_ID 和 INTRANET_ACCESS_KEY_SECRET 环境变量")
        sys.exit(1)

    # 初始化SDK - 使用STS认证方式
    client_config = {
        "access_key_id": access_key_id,
        "access_key_secret": access_key_secret,
    }
    if base_url:
        client_config["base_url"] = base_url

    client = Client(**client_config)

    # 发送消息到Kafka主题
    response = client.connector.send_kafka_message(
        topic="you-topic",
        message={
            "key": "value",
        }
    )
    print(f"错误码: {response.code}")
    print(f"提示信息: {response.msg}")


if __name__ == "__main__":
    main()
