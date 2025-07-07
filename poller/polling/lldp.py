"""MikroTik LLDP polling helper.

Uses the RouterOS API (librouteros) to read LLDP neighbors from the device.
Returns a list of dicts with keys:
    local_mac: MAC address of the local MikroTik interface
    neighbor_mac: MAC address of the connected device (if available)
    neighbor_name: Device ID or name
    port: Port description (optional)

If the MikroTik does not expose LLDP via the API, this module will return an
empty list and log a warning (graceful degradation).
"""
from __future__ import annotations

import os
from typing import List, Dict, Any

from librouteros import connect
from librouteros.exceptions import TrapError

MIKROTIK_HOST = os.getenv("MIKROTIK_HOST", "192.168.1.2")
MIKROTIK_USER = os.getenv("MIKROTIK_USER", "admin")
MIKROTIK_PASS = os.getenv("MIKROTIK_PASS", "password")


def poll_lldp() -> List[Dict[str, Any]]:
    try:
        api = connect(username=MIKROTIK_USER, password=MIKROTIK_PASS, host=MIKROTIK_HOST)
        # RouterOS path for LLDP neighbors
        neighbors = api.path("ip", "neighbor")
        results: list[dict[str, Any]] = []
        for n in neighbors:
            results.append(
                {
                    "local_mac": n.get("mac-address", "").lower(),
                    "neighbor_mac": n.get("neighbor-mac-address", "").lower(),
                    "neighbor_name": n.get("identity", n.get("platform", "")),
                    "port": n.get("interface", ""),
                }
            )
        return results
    except (TrapError, OSError) as exc:
        print(f"[lldp] polling failed: {exc}")
        return []
