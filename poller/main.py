import logging_config  # sets up rotating file + console logging
import threading
import time
import os
from config_loader import load_config
from polling.unifi import poll_unifi
from polling.mikrotik import poll_mikrotik
from polling.lldp import poll_lldp
from processor import process_data
from topology.engine import update_topology
from api.server import start_api

# Load YAML + .env (dotenv handled by docker-compose)
load_config()

POLL_INTERVAL = int(os.getenv("POLL_INTERVAL", "60"))

def poll_loop():
    """Background thread that polls devices and persists data every minute."""
    while True:
        try:
            unifi_data = poll_unifi()
            mikrotik_data = poll_mikrotik()
            lldp_data = poll_lldp()
            combined = unifi_data + mikrotik_data
            # Persist metrics & metadata
            process_data(combined)
            # Update topology snapshot with LLDP
            update_topology(unifi_data, mikrotik_data, lldp_data)
        except Exception as exc:
            print(f"[poller] polling cycle error: {exc}")
        time.sleep(POLL_INTERVAL)

def main():
    # Launch polling loop in background thread
    threading.Thread(target=poll_loop, daemon=True).start()
    # Start FastAPI HTTP server (blocking)
    start_api()

if __name__ == "__main__":
    main()
