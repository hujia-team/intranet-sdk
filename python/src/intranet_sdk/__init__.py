"""MiniEye Intranet SDK - Python client for accessing the Intranet API.

This package provides a simple and efficient way to interact with the MiniEye
Intranet API, supporting STS authentication and various API operations.

Example:
    >>> from intranet_sdk import Client
    >>> client = Client(
    ...     access_key_id="your_access_key_id",
    ...     access_key_secret="your_access_key_secret"
    ... )
    >>> user_info = client.user.get_user_info()
    >>> print(user_info.username)
"""

from .client import Config, HTTPClient
from .intranet import Client
from .models import BaseMsgResp, UserInfo
from .utils import APIError, InternalError, SDKError, set_log_level

__version__ = "0.1.0"

__all__ = [
    "Client",
    "Config",
    "HTTPClient",
    "UserInfo",
    "BaseMsgResp",
    "SDKError",
    "APIError",
    "InternalError",
    "set_log_level",
]
