"""Central logging configuration.

Creates a rotating file handler at /app/logs/poller.log (if writable) and also logs to stdout
so that docker can capture logs. Uses ISO8601 timestamps and includes log level & module.
"""
from __future__ import annotations

import logging
import os
from logging.handlers import RotatingFileHandler
from pathlib import Path

LOG_LEVEL = os.getenv("LOG_LEVEL", "INFO").upper()
LOG_DIR = Path(os.getenv("LOG_DIR", "/app/logs"))
LOG_DIR.mkdir(parents=True, exist_ok=True)
LOG_FILE = LOG_DIR / "poller.log"

FMT = "%(asctime)s %(levelname)s [%(name)s] %(message)s"
DATEFMT = "%Y-%m-%dT%H:%M:%S%z"

root = logging.getLogger()
root.setLevel(LOG_LEVEL)

# Console handler
console = logging.StreamHandler()
console.setFormatter(logging.Formatter(fmt=FMT, datefmt=DATEFMT))
root.addHandler(console)

# Rotating file handler (5 files * 5MB = 25MB max)
try:
    file_handler = RotatingFileHandler(LOG_FILE, maxBytes=5 * 1024 * 1024, backupCount=5)
    file_handler.setFormatter(logging.Formatter(fmt=FMT, datefmt=DATEFMT))
    root.addHandler(file_handler)
except PermissionError:
    root.warning("Unable to write log file at %s; continuing with console-only logging", LOG_FILE)
