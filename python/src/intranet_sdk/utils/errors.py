"""Error types for the Intranet SDK."""

from typing import Optional


class SDKError(Exception):
    """Base exception for all SDK errors."""

    def __init__(self, message: str, cause: Optional[Exception] = None):
        """Initialize SDK error.

        Args:
            message: Error message
            cause: Underlying exception that caused this error
        """
        super().__init__(message)
        self.message = message
        self.cause = cause

    def __str__(self) -> str:
        if self.cause:
            return f"{self.message}: {str(self.cause)}"
        return self.message


class APIError(SDKError):
    """Exception raised when API returns an error."""
    pass


class InternalError(SDKError):
    """Exception raised for internal SDK errors."""
    pass
