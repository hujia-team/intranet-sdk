"""Logging utilities for the SDK."""

import logging
import os
from typing import Any

# Create logger
logger = logging.getLogger("intranet_sdk")

# Set default level from environment or INFO
default_level = os.getenv("INTRANET_SDK_LOG_LEVEL", "INFO").upper()
logger.setLevel(getattr(logging, default_level, logging.INFO))

# Create console handler if not already configured
if not logger.handlers:
    handler = logging.StreamHandler()
    formatter = logging.Formatter(
        "%(asctime)s - %(name)s - %(levelname)s - %(message)s"
    )
    handler.setFormatter(formatter)
    logger.addHandler(handler)


def set_log_level(level: str) -> None:
    """Set the logging level.

    Args:
        level: Logging level (DEBUG, INFO, WARNING, ERROR, CRITICAL)
    """
    logger.setLevel(getattr(logging, level.upper()))


def debug(message: str, *args: Any) -> None:
    """Log a debug message."""
    logger.debug(message, *args)


def info(message: str, *args: Any) -> None:
    """Log an info message."""
    logger.info(message, *args)


def warning(message: str, *args: Any) -> None:
    """Log a warning message."""
    logger.warning(message, *args)


def error(message: str, *args: Any) -> None:
    """Log an error message."""
    logger.error(message, *args)
