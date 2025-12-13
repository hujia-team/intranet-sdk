#!/usr/bin/env python3
"""SDK 测试脚本

快速测试脚本，用于验证 SDK 功能是否正常。

使用方法：
    # 1. 创建 .env 文件
    cp .env.example .env
    # 2. 编辑 .env 填入实际凭证
    # 3. 运行测试
    python test_sdk.py
"""

import os
import sys
from pathlib import Path

# 添加 src 目录到 Python 路径（用于开发测试）
src_path = Path(__file__).parent / "src"
sys.path.insert(0, str(src_path))

# 加载 .env 文件
try:
    from dotenv import load_dotenv

    env_path = Path(__file__).parent / ".env"
    if env_path.exists():
        load_dotenv(env_path)
        print(f"✓ 已加载配置文件: {env_path}")
    else:
        print(f"⚠ 未找到 .env 文件: {env_path}")
        print("请先创建 .env 文件：cp .env.example .env")
        sys.exit(1)
except ImportError:
    print("❌ 需要安装 python-dotenv: pip install python-dotenv")
    sys.exit(1)

from intranet_sdk import Client, set_log_level


def test_user_info():
    """测试获取用户信息"""
    print("\n" + "=" * 60)
    print("测试 1: 获取用户信息")
    print("=" * 60)

    access_key_id = os.getenv("INTRANET_ACCESS_KEY_ID")
    access_key_secret = os.getenv("INTRANET_ACCESS_KEY_SECRET")
    base_url = os.getenv("INTRANET_BASE_URL")

    if not access_key_id or not access_key_secret:
        print("❌ 错误: 未设置 STS 凭证")
        print(
            "请在 .env 文件中配置 INTRANET_ACCESS_KEY_ID 和 INTRANET_ACCESS_KEY_SECRET"
        )
        return False

    try:
        # 创建客户端
        client_config = {
            "access_key_id": access_key_id,
            "access_key_secret": access_key_secret,
        }
        if base_url:
            client_config["base_url"] = base_url

        client = Client(**client_config)

        # 获取用户信息
        user_info = client.user.get_user_info()

        print(f"✓ 用户名       : {user_info.username}")
        print(f"✓ 用户昵称     : {user_info.nickname}")
        print(f"✓ 用户ID       : {user_info.user_id}")
        print(f"✓ 部门名称     : {user_info.department_name}")
        print(f"✓ 角色名称     : {user_info.role_name}")

        return True
    except Exception as e:
        print(f"❌ 测试失败: {e}")
        import traceback

        traceback.print_exc()
        return False


def test_kafka_message():
    """测试发送 Kafka 消息"""
    print("\n" + "=" * 60)
    print("测试 2: 发送 Kafka 消息")
    print("=" * 60)

    access_key_id = os.getenv("INTRANET_ACCESS_KEY_ID")
    access_key_secret = os.getenv("INTRANET_ACCESS_KEY_SECRET")
    base_url = os.getenv("INTRANET_BASE_URL")

    if not access_key_id or not access_key_secret:
        print("❌ 错误: 未设置 STS 凭证")
        return False

    try:
        # 创建客户端
        client_config = {
            "access_key_id": access_key_id,
            "access_key_secret": access_key_secret,
        }
        if base_url:
            client_config["base_url"] = base_url

        client = Client(**client_config)

        # 发送消息
        response = client.connector.send_kafka_message(
            topic="test",
            message={
                "test_key": "test_value",
                "timestamp": "2024-12-13",
            },
        )

        print(f"✓ 错误码       : {response.code}")
        print(f"✓ 提示信息     : {response.msg}")

        if response.is_success():
            print("✓ 消息发送成功")
            return True
        else:
            print(f"❌ 消息发送失败: {response.msg}")
            return False

    except Exception as e:
        print(f"❌ 测试失败: {e}")
        import traceback

        traceback.print_exc()
        return False


def main():
    """运行所有测试"""
    print("\n" + "=" * 60)
    print("Intranet SDK 测试")
    print("=" * 60)

    # 设置日志级别
    log_level = os.getenv("INTRANET_SDK_LOG_LEVEL", "INFO")
    set_log_level(log_level)

    results = []

    # 运行测试
    # results.append(("获取用户信息", test_user_info()))
    results.append(("发送 Kafka 消息", test_kafka_message()))

    # 打印测试结果汇总
    print("\n" + "=" * 60)
    print("测试结果汇总")
    print("=" * 60)

    passed = 0
    failed = 0

    for test_name, result in results:
        status = "✓ 通过" if result else "❌ 失败"
        print(f"{test_name:20s} : {status}")
        if result:
            passed += 1
        else:
            failed += 1

    print("=" * 60)
    print(f"总计: {len(results)} 个测试, {passed} 个通过, {failed} 个失败")
    print("=" * 60)

    # 如果有失败的测试，返回非零退出码
    sys.exit(0 if failed == 0 else 1)


if __name__ == "__main__":
    main()
