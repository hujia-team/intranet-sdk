"""Configuration for the Intranet SDK client."""

from dataclasses import dataclass, field
from typing import Optional


@dataclass
class Config:
    """Configuration for the Intranet API client.

    Attributes:
        base_url: Base URL for API requests. Default: https://intranet.minieye.tech/sys-api
        api_key: API key for authentication (optional)
        access_key_id: Access Key ID for STS authentication (optional)
        access_key_secret: Access Key Secret for STS authentication (optional)
        user_agent: User agent string for requests. Default: minieye-intranet-sdk-python/0.1.0
        timeout: Request timeout in seconds. Default: 30
    """

    base_url: str = "https://intranet.minieye.tech/sys-api"
    api_key: Optional[str] = None
    access_key_id: Optional[str] = None
    access_key_secret: Optional[str] = None
    user_agent: str = "minieye-intranet-sdk-python/0.1.0"
    timeout: int = 30

    def __post_init__(self):
        """Validate configuration after initialization."""
        if not self.base_url:
            raise ValueError("base_url cannot be empty")

        # Remove trailing slash from base_url
        self.base_url = self.base_url.rstrip("/")
