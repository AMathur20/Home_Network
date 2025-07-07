import requests
import os

UNIFI_CONTROLLER = os.getenv("UNIFI_CONTROLLER", "https://192.168.1.1:8443")
USERNAME = os.getenv("UNIFI_USER", "admin")
PASSWORD = os.getenv("UNIFI_PASS", "password")

def login():
    session = requests.Session()
    session.verify = False
    login_data = {
        "username": USERNAME,
        "password": PASSWORD
    }
    resp = session.post(f"{UNIFI_CONTROLLER}/api/login", json=login_data)
    resp.raise_for_status()
    return session

def poll_unifi():
    try:
        session = login()
        clients_resp = session.get(f"{UNIFI_CONTROLLER}/api/s/default/stat/sta")
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
