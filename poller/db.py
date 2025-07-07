"""SQLite database helper for network monitoring metadata.

Creates and migrates the required tables if they do not already exist.
Tables
------
1. devices
   mac TEXT PRIMARY KEY
   hostname TEXT
   first_seen TIMESTAMP
   last_seen TIMESTAMP
   ap_mac TEXT
   switch_mac TEXT
   ip TEXT

2. labels
   mac TEXT PRIMARY KEY REFERENCES devices(mac) ON DELETE CASCADE
   label TEXT NOT NULL

3. topology_snapshots
   id INTEGER PRIMARY KEY AUTOINCREMENT
   ts TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
   graph_json TEXT NOT NULL  -- graph serialized as JSON (nodes, edges, stats)

4. config
   key TEXT PRIMARY KEY
   value TEXT

This module exposes a single `get_connection()` helper that returns a
singleton sqlite3.Connection configured with row factory = sqlite3.Row.
`init_db()` is called upon first import to ensure the schema exists.
"""
from __future__ import annotations

import os
import sqlite3
from pathlib import Path
from typing import Optional

_DB_ENV = os.getenv("SQLITE_PATH", "/app/sqlite/metadata.db")
_DB_PATH = Path(_DB_ENV)

# Ensure parent directory exists when running outside container
_DB_PATH.parent.mkdir(parents=True, exist_ok=True)

_conn: Optional[sqlite3.Connection] = None

SCHEMA = """
PRAGMA foreign_keys = ON;

CREATE TABLE IF NOT EXISTS devices (
    mac TEXT PRIMARY KEY,
    hostname TEXT,
    first_seen TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    last_seen TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    ap_mac TEXT,
    switch_mac TEXT,
    ip TEXT
);

CREATE TABLE IF NOT EXISTS labels (
    mac TEXT PRIMARY KEY REFERENCES devices(mac) ON DELETE CASCADE,
    label TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS topology_snapshots (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    ts TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    graph_json TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS config (
    key TEXT PRIMARY KEY,
    value TEXT
);
"""


def get_connection() -> sqlite3.Connection:
    """Return a singleton connection object with row factory enabled."""
    global _conn
    if _conn is None:
        _conn = sqlite3.connect(_DB_PATH, check_same_thread=False)
        _conn.row_factory = sqlite3.Row
    return _conn


def init_db() -> None:
    """Create tables if they don't exist."""
    conn = get_connection()
    conn.executescript(SCHEMA)
    conn.commit()


# Initialize database on import
init_db()
