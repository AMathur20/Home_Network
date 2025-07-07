"""Load configuration from YAML and inject into environment.

Priority order for configuration values:
1. Existing environment variables (do not override what the user already set)
2. config.yaml (path specified by CONFIG_PATH env or default ./config/config.yaml)
3. Hard-coded defaults in code

Config file structure (example):
---
influx:
  url: "http://influxdb:8086"
  org: "home-net"
  bucket: "network_metrics"
  token: "supersecrettoken"

devices:
  unifi:
    controller: "https://192.168.1.1:8443"
    username: "admin"
    password: "password"
  mikrotik:
    host: "192.168.1.2"
    username: "admin"
    password: "password"

polling:
  interval: 60  # seconds
"""
from __future__ import annotations

import os
from pathlib import Path
from typing import Any

import yaml
from watchdog.observers import Observer
from watchdog.events import FileSystemEventHandler

_DEFAULT_CONFIG_PATH = Path(os.getenv("CONFIG_PATH", "./config/config.yaml"))


def _set_env_if_missing(key: str, value: Any):
    if key not in os.environ and value is not None:
        os.environ[key] = str(value)


def load_config() -> dict[str, Any]:
    if not _DEFAULT_CONFIG_PATH.exists():
        print(f"[config] { _DEFAULT_CONFIG_PATH } not found, using env vars only")
        return {}

    with _DEFAULT_CONFIG_PATH.open("r", encoding="utf-8") as f:
        cfg = yaml.safe_load(f) or {}

    # Influx
    influx = cfg.get("influx", {})
    _set_env_if_missing("INFLUX_URL", influx.get("url"))
    _set_env_if_missing("INFLUX_ORG", influx.get("org"))
    _set_env_if_missing("INFLUX_BUCKET", influx.get("bucket"))
    _set_env_if_missing("INFLUX_TOKEN", influx.get("token"))

    # UniFi
    unifi = cfg.get("devices", {}).get("unifi", {})
    _set_env_if_missing("UNIFI_CONTROLLER", unifi.get("controller"))
    _set_env_if_missing("UNIFI_USER", unifi.get("username"))
    _set_env_if_missing("UNIFI_PASS", unifi.get("password"))

    # MikroTik
    mik = cfg.get("devices", {}).get("mikrotik", {})
    _set_env_if_missing("MIKROTIK_HOST", mik.get("host"))
    _set_env_if_missing("MIKROTIK_USER", mik.get("username"))
    _set_env_if_missing("MIKROTIK_PASS", mik.get("password"))

    # Polling interval
    polling = cfg.get("polling", {})
    _set_env_if_missing("POLL_INTERVAL", polling.get("interval"))

    return cfg


class _ConfigChangeHandler(FileSystemEventHandler):
    def __init__(self):
        super().__init__()

    def on_modified(self, event):
        if Path(event.src_path) == _DEFAULT_CONFIG_PATH:
            print("[config] Detected change in config.yaml, reloading …")
            load_config()


def start_watcher():
    """Start a watchdog observer thread that reloads config on file changes."""
    if not _DEFAULT_CONFIG_PATH.exists():
        return
    obs = Observer()
    obs.daemon = True
    obs.schedule(_ConfigChangeHandler(), str(_DEFAULT_CONFIG_PATH.parent), recursive=False)
    obs.start()

