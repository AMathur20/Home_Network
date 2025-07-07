import os
from pathlib import Path
from poller import config_loader

def test_load_config(tmp_path, monkeypatch):
    yaml_path = tmp_path / "cfg.yaml"
    yaml_path.write_text(
        """
        influx:
          url: http://localhost:8086
          org: test
          bucket: test
          token: t
        polling:
          interval: 30
        devices:
          unifi:
            controller: https://ctrl
            username: u
            password: p
        """
    )
    monkeypatch.setenv("CONFIG_PATH", str(yaml_path))
    monkeypatch.delenv("INFLUX_URL", raising=False)

    cfg = config_loader.load_config()

    assert cfg["polling"]["interval"] == 30
    assert os.getenv("INFLUX_URL") == "http://localhost:8086"
    assert os.getenv("POLL_INTERVAL") == "30"  # set by loader
