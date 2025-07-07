"""Processing and persistence for polled device data.

Responsible for:
• Normalizing records
• Writing bandwidth metrics to InfluxDB
• Upserting device metadata into SQLite
"""
from __future__ import annotations

from datetime import datetime, timezone
from typing import List, Dict, Any

from db import get_connection
from metrics import write_metrics


def _normalize(records: List[Dict[str, Any]]) -> List[Dict[str, Any]]:
    """Ensure required keys are present & consistent."""
    norm: list[dict[str, Any]] = []
    for r in records:
        mac = r.get("mac")
        if not mac:
            # Skip records without MAC – cannot join across sources.
            continue
        norm.append(
            {
                "mac": mac.lower(),
                "hostname": r.get("hostname", ""),
                "ap_mac": r.get("ap_mac", ""),
                "switch_mac": r.get("switch_mac", ""),
                "ip": r.get("ip", ""),
                "rx_bytes": int(r.get("rx_bytes", 0)),
                "tx_bytes": int(r.get("tx_bytes", 0)),
                "interface": r.get("interface", "unknown"),
            }
        )
    return norm


def _upsert_devices(rows: List[Dict[str, Any]]) -> None:
    conn = get_connection()
    cur = conn.cursor()
    now = datetime.now(timezone.utc).isoformat()
    for d in rows:
        cur.execute(
            """
            INSERT INTO devices (mac, hostname, ap_mac, switch_mac, ip, first_seen, last_seen)
            VALUES (:mac, :hostname, :ap_mac, :switch_mac, :ip, :now, :now)
            ON CONFLICT(mac) DO UPDATE SET
                hostname=excluded.hostname,
                ap_mac=excluded.ap_mac,
                switch_mac=excluded.switch_mac,
                ip=excluded.ip,
                last_seen=excluded.last_seen;
            """,
            {**d, "now": now},
        )
    conn.commit()


def _prepare_metric_points(rows: List[Dict[str, Any]]):
    ts = datetime.now(timezone.utc)
    for d in rows:
        yield {
            "measurement": "bandwidth",
            "time": ts,
            "tags": {
                "mac": d["mac"],
                "interface": d["interface"],
            },
            "fields": {
                "rx_bytes": d["rx_bytes"],
                "tx_bytes": d["tx_bytes"],
            },
        }


def process_data(devices: List[Dict[str, Any]]) -> None:
    """Normalize, persist metadata, and push metrics to InfluxDB."""
    normalized = _normalize(devices)
    if not normalized:
        return

    # Persist metadata
    _upsert_devices(normalized)

    # Push metrics to InfluxDB
    write_metrics(_prepare_metric_points(normalized))
