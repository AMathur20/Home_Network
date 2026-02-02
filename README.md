# Home Network Monitor (HNM) v2

**Home Network Monitor (HNM)** is a high-performance, containerized observability suite designed for multi-vendor home labs. It provides a "Dark NOC" dashboard to visualize network topology, link health, and real-time bandwidth across MikroTik, UniFi, and EdgeRouter hardware.

## 🚀 Features

- **"Dark NOC" Dashboard**: Premium React-based UI with interactive D3 topology maps.
- **Link Classification**: Automatic detection of 10G, 1G, and Wireless links based on SNMP data.
- **Mobile Optimization**: Dedicated **Priority List** highlighting high-bandwidth and "Down" links for small screens.
- **Auto-Discovery**: Automatic topology generation using SNMP LLDP-MIB.
- **Precision Polling**: Delta-based throughput calculation for accurate bits-per-second reporting.
- **Portainer Ready**: Easy deployment with Docker Compose and integrated configuration editor.

## 🛠 Technology Stack

- **Backend**: Go (Golang)
- **Frontend**: React + Tailwind CSS + D3.js
- **Database**: DuckDB
- **Config**: YAML (`config.yaml`, `topology.yaml`)
- **Container**: Docker + Docker Compose

## 📖 Documentation

- [Product Requirements Document (PRD)](prd.md)

---

## 📦 Deployment on Portainer

HNM v2 is designed to be deployed as a "Stack" in Portainer for easy management.

### 1. Preparation
Ensure you have a directory on your host to store configurations (e.g., `./hnm/config`).

### 2. Portainer Stack Setup
1. Log in to your Portainer instance.
2. Go to **Stacks** > **Add stack**.
3. Name your stack (e.g., `hnm`).
4. Select **Web editor** and paste the following configuration:

```yaml
services:
  hnm-core:
    build: . # Or use image: ghcr.io/[user]/hnm-core:latest
    container_name: hnm-core
    ports:
      - "8080:8080"
    volumes:
      - /path/to/your/config:/app/config
      - hnm-data:/app/data
    environment:
      - HNM_CONFIG_PATH=/app/config/config.yaml
      - HNM_DB_PATH=/app/data/hnm.db
    restart: always

  hnm-editor:
    image: filebrowser/filebrowser:latest
    container_name: hnm-editor
    ports:
      - "8081:80"
    volumes:
      - /path/to/your/config:/srv
    environment:
      - FB_BASEURL=/config
    restart: always

volumes:
  hnm-data:
```

> [!IMPORTANT]
> Change `/path/to/your/config` to the absolute path on your host where you want to store `config.yaml` and `topology.yaml`. HNM will create a `hnm.db` in the named volume `hnm-data`.

### 3. Deployment
Click **Deploy the stack**. Portainer will build (or pull) the images and start the containers.

---

## ⚙️ Configuration

### Initial Setup
1. Access the **Config Editor** at `http://[YOUR-IP]:8081`.
2. Create or upload your `config.yaml` in the root of the file browser.
3. Example `config.yaml`:
   ```yaml
   poller:
     live: 5
     history: 60
   devices:
     - name: "core-router"
       host: "192.168.1.1"
       type: "mikrotik"
       snmp:
         version: "v2c"
         community: "public"
         port: 161
   ```
4. Restart the `hnm-core` container in Portainer to apply the changes.

### Topology
HNM will automatically attempt to crawl your network via LLDP on the first boot. You can manually override or correct links by editing `topology.yaml` via the Config Editor. The dashboard will hot-reload automatically when you save changes.

## 🤝 License
MIT
