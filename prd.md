# PRD: Home Network Monitor (HNM) v2

## 1. Project Vision

**Home Network Monitor (HNM)** is a high-performance, containerized observability suite designed for multi-vendor home labs. It provides a "Dark NOC" dashboard to visualize network topology, link health, and real-time/historical bandwidth across MikroTik, UniFi, and EdgeRouter hardware.

---

## 2. Core Functional Requirements

### 2.1 Polling Engine (The Core)

* **Protocol:** SNMP (v2c/v3) and Vendor APIs (RouterOS API).
* **Live View (5s):** Concurrent polling of interface counters for real-time throughput.
* **History (60s):** Data persistence for time-series analysis.
* **Link Detection:** Immediate status reporting (Up/Down) based on administrative and operational states.

### 2.2 Intelligent Topology Logic

* **Auto-Discovery:** One-time crawl using **LLDP-MIB** to identify neighbor relationships.
* **Topology Persistence:** Findings are written to `topology.yaml`.
* **Static Overrides:** The Go backend reads `topology.yaml` on boot. Manual corrections in the YAML take precedence over auto-discovery logic.
* **Link Classification:** Automatic detection of link types based on interface names (e.g., `sfp-sfpplus` = 10G) with YAML-based manual fallback.

### 2.3 The "Dark NOC" Dashboard

* **Visual Language:**
* **10Gb SFP+:** Thick Neon Cyan, pulsing "flow" animation.
* **1Gb SFP:** Medium Emerald Green, steady flow.
* **1Gb Ethernet:** Medium Forest Green, static/slow flow.
* **Wireless:** Dotted Orange/Yellow, static.


* **Dual-View Engine:** Toggle between **Force-Directed Graph** (dynamic) and **Grid Layout** (structured).
* **Mobile Responsiveness:** Auto-collapse map to a prioritized list of high-bandwidth or "Down" links on small screens.

---

## 3. Technical Specifications

| Component | Technology |
| --- | --- |
| **Backend** | **Go (Golang)** - chosen for high-concurrency polling and low memory footprint. |
| **Frontend** | **React + Tailwind CSS** (Dark Mode focus). |
| **Database** | **DuckDB** - Single-file, high-performance time-series storage. |
| **Config** | YAML (`config.yaml` for auth, `topology.yaml` for the map). |
| **Container** | Docker (Distroless or Alpine base for security/speed). |

---

## 4. Configuration & Deployment

### 4.1 Sidecar Configuration Management

To ensure the YAML files are easily editable within Portainer, the stack includes a **FileBrowser** sidecar. This allows for in-browser editing of the network structure without needing SSH or CLI access.

### 4.2 Portainer Deployment Template (Draft)

> **Note to Deployer:** Port mapping (e.g., 80 and 8081) should be adjusted to fit your specific network environment and avoid conflicts.

```yaml
services:
  hnm-core:
    image: ghcr.io/[user]/hnm-core:latest
    container_name: hnm-core
    ports:
      - "80:8080" # CHANGE THIS: HostPort:ContainerPort
    volumes:
      - ./config:/app/config
    restart: always

  hnm-editor:
    image: filebrowser/filebrowser:latest
    container_name: hnm-editor
    ports:
      - "8081:80" # CHANGE THIS: HostPort:ContainerPort
    volumes:
      - ./config:/srv
    environment:
      - FB_BASEURL=/config
    restart: always

```

---

## 5. Phased Roadmap for "Antigravity"

### Phase 1: Go Polling Foundation

* Build the Go SNMP/RouterOS poller.
* Implement `config.yaml` parsing.
* Test throughput calculation ().

### Phase 2: Topology Generation

* Implement LLDP crawler logic.
* Generate the initial `topology.yaml`.
* Build the `fsnotify` file watcher to hot-reload the map when the YAML is edited.

### Phase 3: Dashboard UI

* Create the React frontend with the "Dark NOC" theme.
* Implement the link styling logic (Cyan/Green/Orange).
* Add the Force-Directed vs. Grid toggle.

### Phase 4: Release & Packaging

* Finalize Docker multi-arch builds (ARM64 support for Raspberry Pi).
* Create GitHub Release workflow for versioned images.

