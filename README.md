# Home Network Monitor (HNM) v2

**Home Network Monitor (HNM)** is a high-performance, containerized observability suite designed for multi-vendor home labs. It provides a "Dark NOC" dashboard to visualize network topology, link health, and real-time bandwidth across MikroTik, UniFi, and EdgeRouter hardware.

## 🚀 Features

- **"Dark NOC" Dashboard**: Premium React-based UI with interactive D3 topology maps.
- **Link Classification**: Automatic detection of 10G, 1G, and Wireless links based on SNMP data.
- **Mobile Optimization**: Dedicated **Priority List** highlighting high-bandwidth and "Down" links for small screens.
- **Auto-Discovery**: Automatic topology generation using SNMP LLDP-MIB.
- **Precision Polling**: Delta-based throughput calculation for accurate bits-per-second reporting.
- **Docker Native**: Built for seamless deployment on Ubuntu and Linux servers.

## 🛠 Technology Stack

- **Backend**: Go (Golang)
- **Frontend**: React + Tailwind CSS + D3.js
- **Database**: DuckDB
- **Config**: YAML (`config.yaml`, `topology.yaml`)
- **Container**: Docker + Docker Compose

## 📖 Documentation

- [Product Requirements Document (PRD)](prd.md)

---

## 📦 Deployment Guide (Ubuntu / CLI)

Deploying HNM v2 on an Ubuntu server via the command line is the recommended approach for maximum stability and performance.

### 1. Prerequisites (Ubuntu)
Ensure Docker and Docker Compose are installed:
```bash
sudo apt update
sudo apt install docker.io docker-compose -y
sudo usermod -aG docker $USER
# Log out and log back in for group changes to take effect
```

### 2. Prepare Host Storage
Create a dedicated folder on your host to store configurations and data.
```bash
# Create config directory
mkdir -p ~/hnm/config

# Set ownership (ensures Docker can write metrics/logs)
sudo chown -R $USER:$USER ~/hnm
```

### 3. Configure the Environment
HNM uses a `.env` file to manage paths without hardcoding them into the configuration. Create a file named `.env` in the project root:
```bash
echo "HNM_CONFIG_DIR=$HOME/hnm/config" > .env
```

### 4. Build and Launch
Run the following command from the project root:
```bash
# Build the images and start services in detached mode
docker-compose up -d --build
```

### 5. Access the Dashboard
- **HNM Dashboard**: `http://<your-server-ip>:8080`
- **Config Editor**: `http://<your-server-ip>:8081`

---

## ⚙️ Configuration

### Initial Setup
1. Open the **Config Editor** at `http://<server-ip>:8081`.
2. Create or upload your `config.yaml`.
3. **Example `config.yaml`**:
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
4. Restart the poller: `docker-compose restart hnm-core`

### Topology
HNM will automatically attempt to crawl your network via LLDP on the first boot. You can manually override or correct links by editing `topology.yaml` via the Config Editor.

---

## ⛵ Alternative: Deployment on Portainer

If you prefer using Portainer:
1. Go to **Stacks** > **Add stack**.
2. Select **Repository** (point to your Git URL).
3. In **Environment variables**, add:
   - `HNM_CONFIG_DIR`: `/your/host/path/to/config`
4. Click **Deploy the stack**.

## 🤝 License

Apache 2.0