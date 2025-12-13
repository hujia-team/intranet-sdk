"""STS (Security Token Service) authentication utilities."""

import hashlib
from datetime import datetime


def md5sum(content: str) -> str:
    """Calculate MD5 hash of content.

    Args:
        content: Content to hash

    Returns:
        MD5 hash as 32-character hex string
    """
    hasher = hashlib.md5()
    if not isinstance(content, bytes):
        content = content.encode("utf-8")
    hasher.update(content)
    return hasher.hexdigest()


def get_sts_token(access_key_id: str, access_key_secret: str) -> str:
    """Generate STS token for authentication.

    The token is generated using MD5 hash of the format:
    {year}-{month}-{day}_{access_key_id}_{access_key_secret}_{hour}

    Args:
        access_key_id: Access Key ID
        access_key_secret: Access Key Secret

    Returns:
        32-character MD5 hash token

    Raises:
        ValueError: If system timezone is not UTC+8
    """
    # Get current time with timezone
    now = datetime.now().astimezone()

    # Check timezone (must be UTC+8)
    timezone_offset = now.utcoffset().total_seconds() / 3600
    if timezone_offset != 8:
        raise ValueError(
            f"系统时区必须为东八区(UTC+8)，当前时区偏移为{timezone_offset}小时"
        )

    # Format: YYYY-M-D (without leading zeros)
    today = f"{now.year}-{now.month}-{now.day}"

    # Generate token string
    txt = f"{today}_{access_key_id}_{access_key_secret}_{now.hour}"

    # Return MD5 hash
    return md5sum(txt)
