"""Topology Engine

Creates a simple graph model (nodes & edges) from polled UniFi and
MikroTik data and stores a JSON snapshot in the SQLite DB. The schema is
created by `poller.db` and contains:
    topology_snapshots(id, ts, graph_json)

The graph structure generated is intentionally minimal – suitable for the
Grafana Node-Graph plugin:
    {
        "nodes": [{"id": "<mac>", "label": "<hostname/iface>"}, ...],
        "edges": [{"source": "<mac>", "target": "<mac>"}, ...]
    }

In the future this can be enhanced with interface speeds, port numbers,
latency metrics, etc.
"""
from __future__ import annotations

import json
from typing import List, Dict, Any, Set

from db import get_connection
from metrics import write_metrics


def _add_node(nodes: List[Dict[str, str]], seen: Set[str], mac: str, label: str = "") -> None:
    mac = mac.lower()
    if mac and mac not in seen:
        nodes.append({"id": mac, "label": label or mac})
        seen.add(mac)


def build_graph(unifi_data: List[Dict[str, Any]], mikrotik_data: List[Dict[str, Any]], lldp_data: List[Dict[str, Any]] | None = None):
    """Return a tuple (nodes, edges) representing the current topology."""
    nodes: list[dict[str, str]] = []
    edges: list[dict[str, str]] = []
    seen_nodes: set[str] = set()

    # Map of AP MAC to switch/parent if known (could be enriched later)

    # UniFi clients → AP edges
    for d in unifi_data:
        client_mac = d.get("mac", "").lower()
        ap_mac = d.get("ap_mac", "").lower()
        if not client_mac:
            continue
        _add_node(nodes, seen_nodes, client_mac, d.get("hostname", client_mac))
        if ap_mac:
            _add_node(nodes, seen_nodes, ap_mac, ap_mac)
            edges.append({"source": client_mac, "target": ap_mac})

    # Mikrotik interfaces – treat each interface as node, attach to host router MAC placeholder
    for i in mikrotik_data:
        iface_mac = i.get("mac", "").lower()
        if not iface_mac:
            continue
        _add_node(nodes, seen_nodes, iface_mac, i.get("name", iface_mac))
        # simplistic: connect interface to a dummy router node "mikrotik"
        _add_node(nodes, seen_nodes, "mikrotik", "MikroTik")
        edges.append({"source": iface_mac, "target": "mikrotik"})

    # LLDP wired edges (switch/router connections)
    if lldp_data:
        for l in lldp_data:
            local = l.get("local_mac", "").lower()
            neighbor = l.get("neighbor_mac", "").lower()
            if local:
                _add_node(nodes, seen_nodes, local, local)
            if neighbor:
                _add_node(nodes, seen_nodes, neighbor, l.get("neighbor_name", neighbor))
            if local and neighbor:
                edges.append({"source": local, "target": neighbor})

    return nodes, edges


def update_topology(unifi_data: List[Dict[str, Any]], mikrotik_data: List[Dict[str, Any]], lldp_data: List[Dict[str, Any]] | None = None):
    """Build graph from latest polled data and insert snapshot into SQLite."""
    nodes, edges = build_graph(unifi_data, mikrotik_data, lldp_data)
    graph_json = json.dumps({"nodes": nodes, "edges": edges})
    # Persist to SQLite
    conn = get_connection()
    conn.execute(
        "INSERT INTO topology_snapshots (graph_json) VALUES (?);",
        (graph_json,),
    )
    conn.commit()

    # Publish to InfluxDB so Grafana NodeGraph can query
    write_metrics([
        {
            "measurement": "topology",
            "tags": {},
            "fields": {"graph": graph_json},
        }
    ])

    return nodes, edges

