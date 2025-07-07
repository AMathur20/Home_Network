from librouteros import connect
import os

MIKROTIK_HOST = os.getenv("MIKROTIK_HOST", "192.168.1.2")
MIKROTIK_USER = os.getenv("MIKROTIK_USER", "admin")
MIKROTIK_PASS = os.getenv("MIKROTIK_PASS", "password")

def poll_mikrotik():
    try:
        api = connect(username=MIKROTIK_USER, password=MIKROTIK_PASS, host=MIKROTIK_HOST)
        interfaces = api.path("interface")
        return [{
            "name": iface.get("name"),
            "mac": iface.get("mac-address"),
            "rx_bytes": int(iface.get("rx-byte", 0)),
            "tx_bytes": int(iface.get("tx-byte", 0)),
            "interface": "wired"
        } for iface in interfaces]
    except Exception as e:
        print(f"[mikrotik] polling failed: {e}")
        return []
