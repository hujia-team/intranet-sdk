"""User-related data models."""

from dataclasses import dataclass
from typing import Optional


@dataclass
class UserInfo:
    """User information.

    Attributes:
        user_id: User unique identifier (UUID)
        username: Username
        nickname: Nickname
        avatar: Avatar URL
        home_path: Home directory path
        role_name: Role name
        department_name: Department name
    """

    user_id: Optional[str] = None
    username: Optional[str] = None
    nickname: Optional[str] = None
    avatar: Optional[str] = None
    home_path: Optional[str] = None
    role_name: Optional[str] = None
    department_name: Optional[str] = None

    @classmethod
    def from_dict(cls, data: dict) -> "UserInfo":
        """Create UserInfo instance from dictionary.

        Args:
            data: Dictionary containing user data

        Returns:
            UserInfo instance
        """
        return cls(
            user_id=data.get("userId"),
            username=data.get("username"),
            nickname=data.get("nickname"),
            avatar=data.get("avatar"),
            home_path=data.get("homePath"),
            role_name=data.get("roleName"),
            department_name=data.get("departmentName"),
        )

    def get_username(self) -> str:
        """Get username.

        Returns:
            Username or empty string if not set
        """
        return self.username or ""

    def get_nickname(self) -> str:
        """Get nickname.

        Returns:
            Nickname or empty string if not set
        """
        return self.nickname or ""
