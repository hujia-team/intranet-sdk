"""HTTP client for making API requests."""

import json
from typing import Any, Dict, Optional
from urllib.parse import urljoin

import requests

from ..utils import APIError, InternalError, debug, error
from ..utils.sts import get_sts_token
from .config import Config


class HTTPClient:
    """HTTP client for Intranet API requests.

    This class handles all HTTP communication with the API, including
    authentication, request signing, and error handling.
    """

    def __init__(self, config: Config):
        """Initialize HTTP client with configuration.

        Args:
            config: Client configuration
        """
        self.config = config
        self.session = requests.Session()
        self.session.headers.update(
            {
                "User-Agent": config.user_agent,
                "Content-Type": "application/json",
                "Accept": "application/json",
            }
        )

    def _build_url(self, path: str) -> str:
        """Build full URL from path.

        Args:
            path: API endpoint path

        Returns:
            Full URL
        """
        # Ensure path starts with /
        if not path.startswith("/"):
            path = f"/{path}"

        # Ensure base_url ends with / for proper joining
        base = self.config.base_url
        if not base.endswith("/"):
            base += "/"

        # Remove leading / from path to avoid double slashes
        if path.startswith("/"):
            path = path[1:]

        # Simple string concatenation to avoid urljoin replacing base path
        return base + path

    def _add_auth_headers(self) -> Dict[str, str]:
        """Add authentication headers to request.

        Returns:
            Dictionary of authentication headers
        """
        headers = {}

        # Priority: API Key > STS authentication
        if self.config.api_key:
            headers["Authorization"] = f"Bearer {self.config.api_key}"
        elif self.config.access_key_id and self.config.access_key_secret:
            # Use STS authentication
            headers["x-sts-uid"] = self.config.access_key_id
            headers["x-sts-token"] = get_sts_token(
                self.config.access_key_id,
                self.config.access_key_secret
            )
            debug(f"STS auth - UID: {self.config.access_key_id}")

        return headers

    def _do_request(
        self,
        method: str,
        path: str,
        data: Optional[Dict[str, Any]] = None,
        params: Optional[Dict[str, Any]] = None,
    ) -> Any:
        """Perform HTTP request.

        Args:
            method: HTTP method (GET, POST, etc.)
            path: API endpoint path
            data: Request body data (for POST/PUT)
            params: Query parameters

        Returns:
            Response data

        Raises:
            APIError: If API returns an error
            InternalError: If request fails
        """
        url = self._build_url(path)
        headers = self._add_auth_headers()

        debug(f"{method} {url}")

        try:
            response = self.session.request(
                method=method,
                url=url,
                json=data,
                params=params,
                headers=headers,
                timeout=self.config.timeout,
            )
            response.raise_for_status()
            return response.json()
        except requests.exceptions.HTTPError as e:
            error(f"HTTP error: {e}")
            raise InternalError(f"HTTP request failed: {path}", e)
        except requests.exceptions.RequestException as e:
            error(f"Request failed: {e}")
            raise InternalError(f"Request failed: {path}", e)
        except json.JSONDecodeError as e:
            error(f"Failed to decode JSON response: {e}")
            raise InternalError("Invalid JSON response", e)

    def get(
        self,
        path: str,
        params: Optional[Dict[str, Any]] = None,
    ) -> Any:
        """Make a GET request.

        Args:
            path: API endpoint path
            params: Query parameters

        Returns:
            Response data

        Raises:
            APIError: If API returns an error
            InternalError: If request fails
        """
        return self._do_request("GET", path, params=params)

    def post(
        self,
        path: str,
        data: Optional[Dict[str, Any]] = None,
        params: Optional[Dict[str, Any]] = None,
    ) -> Any:
        """Make a POST request.

        Args:
            path: API endpoint path
            data: Request body data
            params: Query parameters

        Returns:
            Response data

        Raises:
            APIError: If API returns an error
            InternalError: If request fails
        """
        return self._do_request("POST", path, data=data, params=params)

    def put(
        self,
        path: str,
        data: Optional[Dict[str, Any]] = None,
        params: Optional[Dict[str, Any]] = None,
    ) -> Any:
        """Make a PUT request.

        Args:
            path: API endpoint path
            data: Request body data
            params: Query parameters

        Returns:
            Response data

        Raises:
            APIError: If API returns an error
            InternalError: If request fails
        """
        return self._do_request("PUT", path, data=data, params=params)

    def delete(
        self,
        path: str,
        params: Optional[Dict[str, Any]] = None,
    ) -> Any:
        """Make a DELETE request.

        Args:
            path: API endpoint path
            params: Query parameters

        Returns:
            Response data

        Raises:
            APIError: If API returns an error
            InternalError: If request fails
        """
        return self._do_request("DELETE", path, params=params)
