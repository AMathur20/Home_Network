import os
import requests





def _env(key: str, default: str | None = None) -> str:
    """Read env at call-time so changes made by config_loader take effect."""
    return os.getenv(key, default) or default

def _controller() -> str:
    """Return UniFi controller base URL."""
    return _env("UNIFI_CONTROLLER", "https://192.168.1.1:8443")

def login():
    controller = _controller()
    username = _env("UNIFI_USER", "admin")
    password = _env("UNIFI_PASS", "password")

    session = requests.Session()
    session.verify = False
    login_data = {
        "username": username,
        "password": password
    }
    resp = session.post(f"{controller}/api/login", json=login_data)
    resp.raise_for_status()
    return session

def poll_unifi():
    controller = _controller()
    try:
        session = login()
        clients_resp = session.get(f"{controller}/api/s/default/stat/sta")
        clients = clients_resp.json().get("data", [])
        return [{
            "mac": c["mac"],
            "hostname": c.get("hostname", ""),
            "ap_mac": c.get("ap_mac", ""),
            "rx_bytes": c.get("rx_bytes", 0),
            "tx_bytes": c.get("tx_bytes", 0),
            "interface": "wireless"
        } for c in clients]
    except Exception as e:
        print(f"[unifi] polling failed: {e}")
        return []
