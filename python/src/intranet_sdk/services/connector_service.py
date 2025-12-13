"""Connector service for messaging operations."""

import json
from typing import TYPE_CHECKING, Any

from ..models import BaseMsgResp
from ..utils import APIError, InternalError, debug, error

if TYPE_CHECKING:
    from ..client import HTTPClient


class ConnectorService:
    """Service for connector-related operations."""

    def __init__(self, http_client: "HTTPClient"):
        """Initialize connector service.

        Args:
            http_client: HTTP client for making requests
        """
        self.http_client = http_client

    def send_kafka_message(self, topic: str, message: Any) -> BaseMsgResp:
        """Send a message to Kafka topic.

        Args:
            topic: Kafka topic name
            message: Message to send (will be JSON serialized)

        Returns:
            BaseMsgResp with response status

        Raises:
            APIError: If API returns an error
            InternalError: If request fails or JSON serialization fails
        """
        debug(f"Sending message to Kafka topic: {topic}")

        # Serialize message to JSON
        try:
            message_json = json.dumps(message)
        except (TypeError, ValueError) as e:
            error(f"JSON serialization failed: {e}")
            raise InternalError("Failed to serialize message", e)

        # Prepare request data
        data = {
            "topic": topic,
            "message": message_json,
        }

        # Make API request
        response = self.http_client.post("/connector/kafka/send-topic-message", data=data)

        if response.get("code") != 0:
            msg = response.get("msg", "Unknown error")
            error(f"API error: {msg}")
            return BaseMsgResp(code=response.get("code", -1), msg=msg)

        debug(f"Sent message to Kafka topic: {topic} successfully")
        return BaseMsgResp(code=response.get("code", 0), msg=response.get("msg", "Success"))
