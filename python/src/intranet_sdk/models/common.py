"""Common data models."""

from dataclasses import dataclass
from typing import Optional


@dataclass
class BaseMsgResp:
    """Base response without data.

    Attributes:
        code: Error code (0 for success)
        msg: Response message
    """

    code: int
    msg: str

    @classmethod
    def from_dict(cls, data: dict) -> "BaseMsgResp":
        """Create instance from dictionary.

        Args:
            data: Dictionary containing response data

        Returns:
            BaseMsgResp instance
        """
        return cls(
            code=data.get("code", 0),
            msg=data.get("msg", ""),
        )

    def is_success(self) -> bool:
        """Check if response indicates success.

        Returns:
            True if code is 0, False otherwise
        """
        return self.code == 0
