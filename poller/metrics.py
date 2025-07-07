"""Utility functions for writing metrics to InfluxDB 2.x.

Relies on environment variables to configure the connection:
INFLUX_URL – http(s)://host:port (default: http://influxdb:8086)
INFLUX_TOKEN – API token (default: supersecrettoken)
INFLUX_ORG   – organisation name (default: home-net)
INFLUX_BUCKET – bucket (default: network_metrics)
"""
from __future__ import annotations

import os
from typing import Iterable, Mapping
from datetime import datetime, timezone

from influxdb_client import InfluxDBClient, Point, WriteOptions

_INFLUX_URL = os.getenv("INFLUX_URL", "http://influxdb:8086")
_INFLUX_TOKEN = os.getenv("INFLUX_TOKEN", "supersecrettoken")
_INFLUX_ORG = os.getenv("INFLUX_ORG", "home-net")
_INFLUX_BUCKET = os.getenv("INFLUX_BUCKET", "network_metrics")

_client: InfluxDBClient | None = None


def _get_client() -> InfluxDBClient:
    global _client
    if _client is None:
        _client = InfluxDBClient(url=_INFLUX_URL, token=_INFLUX_TOKEN, org=_INFLUX_ORG)
    return _client


def write_metrics(metrics: Iterable[Mapping[str, object]]) -> None:
    """Write a sequence of metric dicts to InfluxDB.

    Each metric mapping must contain at minimum:
        measurement: str
        tags: dict[str, str]
        fields: dict[str, int|float|str]
        time: datetime | None – optional; defaults to now() UTC.
    """
    client = _get_client()
    write_api = client.write_api(write_options=WriteOptions(batch_size=500, flush_interval=1000))

    points: list[Point] = []
    for m in metrics:
        measurement = m["measurement"]
        tags = m.get("tags", {})
        fields = m.get("fields", {})
        ts: datetime | None = m.get("time")
        if ts is None:
            ts = datetime.now(timezone.utc)
        pt = Point(measurement).time(ts)
        for k, v in tags.items():
            pt = pt.tag(k, str(v))
        for k, v in fields.items():
            pt = pt.field(k, v)
        points.append(pt)

    if points:
        write_api.write(bucket=_INFLUX_BUCKET, record=points)
