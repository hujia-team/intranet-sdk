"""Main client for the Intranet SDK."""

from typing import Optional

from .client import Config, HTTPClient
from .services import ConnectorService, UserService


class Client:
    """Client for the MiniEye Intranet API.

    This is the main entry point for interacting with the API. It provides
    access to various services through its properties.

    Attributes:
        user: User service for user-related operations
        connector: Connector service for messaging operations

    Example:
        >>> client = Client(
        ...     access_key_id="your_key_id",
        ...     access_key_secret="your_key_secret"
        ... )
        >>> user_info = client.user.get_user_info()
        >>> print(user_info.username)
    """

    def __init__(
        self,
        base_url: str = "https://intranet.minieye.tech/sys-api",
        api_key: Optional[str] = None,
        access_key_id: Optional[str] = None,
        access_key_secret: Optional[str] = None,
        user_agent: str = "minieye-intranet-sdk-python/0.1.0",
        timeout: int = 30,
    ):
        """Initialize the Intranet API client.

        Args:
            base_url: Base URL for API requests. Default: https://intranet.minieye.tech/sys-api
            api_key: API key for authentication (optional)
            access_key_id: Access Key ID for STS authentication (optional)
            access_key_secret: Access Key Secret for STS authentication (optional)
            user_agent: User agent string. Default: minieye-intranet-sdk-python/0.1.0
            timeout: Request timeout in seconds. Default: 30

        Raises:
            ValueError: If configuration is invalid
        """
        # Create configuration
        config = Config(
            base_url=base_url,
            api_key=api_key,
            access_key_id=access_key_id,
            access_key_secret=access_key_secret,
            user_agent=user_agent,
            timeout=timeout,
        )

        # Create HTTP client
        self._http_client = HTTPClient(config)

        # Initialize services
        self.user = UserService(self._http_client)
        self.connector = ConnectorService(self._http_client)
