"""User service for managing user-related operations."""

from typing import TYPE_CHECKING

from ..models import UserInfo
from ..utils import APIError, debug, error

if TYPE_CHECKING:
    from ..client import HTTPClient


class UserService:
    """Service for user-related operations."""

    def __init__(self, http_client: "HTTPClient"):
        """Initialize user service.

        Args:
            http_client: HTTP client for making requests
        """
        self.http_client = http_client

    def get_user_info(self) -> UserInfo:
        """Get current user's information.

        Returns:
            UserInfo object containing user data

        Raises:
            APIError: If API returns an error
            InternalError: If request fails
        """
        debug("Getting current user info")

        response = self.http_client.get("/user/info")

        if response.get("code") != 0:
            msg = response.get("msg", "Unknown error")
            error(f"API error: {msg}")
            raise APIError(msg)

        debug("Got user info successfully")
        return UserInfo.from_dict(response.get("data", {}))
