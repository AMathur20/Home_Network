from poller.topology.engine import build_graph

def test_build_graph_simple():
    unifi_data = [
        {"mac": "aa:bb", "ap_mac": "11:22", "hostname": "device1"},
        {"mac": "cc:dd", "ap_mac": "11:22", "hostname": "device2"},
    ]
    mik_data = [
        {"mac": "11:22", "hostname": "ap1"},
    ]
    lldp_data = [
        {"local_mac": "11:22", "neighbor_mac": "ee:ff", "local_port": "1"}
    ]

    nodes, edges = build_graph(unifi_data, mik_data, lldp_data)
    macs = {n["id"] for n in nodes}
    assert {"aa:bb", "cc:dd", "11:22", "ee:ff"}.issubset(macs)
    # Expect at least three edges (two wireless, one wired)
    assert len(edges) >= 3
