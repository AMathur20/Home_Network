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

### 1. Preparation: Setting Up Host Storage

Before deploying, you must create a directory on your Docker host to store the YAML configuration files. This ensures your settings (like device IPs and network topology) are persistent and reachable by both the HNM core and the FileBrowser sidecar.

#### On Linux or Mac:
1. Open your terminal on the Docker host.
2. Create a configuration directory (e.g., in your home folder):
   ```bash
   mkdir -p ~/hnm/config
   ```
3. Set appropriate permissions (ensures the Docker container can write to this folder):
   ```bash
   chmod -R 755 ~/hnm/config
   ```
4. **Note the Absolute Path**: Run `pwd` inside that folder. You will need this path (e.g., `/home/user/hnm/config`) for the Portainer stack configuration.

#### Manual File Creation (Optional):
The integrated **FileBrowser** will allow you to create these files via the UI, but you can also pre-create an empty `config.yaml` to be ready:
```bash
touch ~/hnm/config/config.yaml
```

### 2. Portainer Stack Setup

> [!WARNING]
> **Why did the build fail?** If you are using the Portainer **Web Editor**, the `build: .` command will fail because Portainer does not have access to your local source code. You have two options:

#### Option A: Build Locally (Easiest for most)
Run this command in your project root on your terminal **before** deploying in Portainer:
```bash
docker build -t hnm-core:latest .
```
Once the build is finished, you can use `image: hnm-core:latest` in your Portainer stack, and it will find the image you just built.

#### Option B: Deploy via Git
Instead of "Web editor", select **Repository** in Portainer and point it to your GitHub/GitLab URL. 

**How to handle paths in Git mode:**
- Portainer will clone the repo and use the `docker-compose.yml` inside it.
- **Critical**: You still need to ensure the `volumes` in the `docker-compose.yml` point to your actual host paths. 
- You can either:
    1.  **Edit the file in your repo** to include your specific paths before pushing.
    2.  **Use Portainer Environment Variables**: Use variables like `${HNM_CONFIG_DIR}` in the YAML and define them in the Portainer "Environment variables" section when creating the stack.

### 3. Deployment
1. Log in to your Portainer instance.
2. Go to **Stacks** > **Add stack**.
3. Name your stack (e.g., `hnm`).
4. Paste the following configuration (if using **Option A**):

```yaml
services:
  hnm-core:
    image: hnm-core:latest 
    container_name: hnm-core
    ports:
      - "8080:8080"
    volumes:
      - ${HNM_CONFIG_DIR}:/app/config
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
      - ${HNM_CONFIG_DIR}:/srv
    environment:
      - FB_BASEURL=/config
    restart: always

volumes:
  hnm-data:
```

> [!IMPORTANT]
> **Environment Variables**: When creating the stack in Portainer, you **must** define a variable named `HNM_CONFIG_DIR` and set it to your absolute host path (e.g., `/home/user/hnm/config`). This ensures your configuration is persistent and accessible to both services.

### 4. Click 'Deploy'
Click **Deploy the stack**.

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
