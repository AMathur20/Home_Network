import sqlite3
from poller import db


def test_schema_tables_exist(tmp_path):
    # Use a temporary DB path
    db_path = tmp_path / "test.db"
    original_env = db.os.environ.get("SQLITE_PATH")
    db.os.environ["SQLITE_PATH"] = str(db_path)

    # Force reconnect
    db._conn = None
    conn = db.get_connection()

    # Assert each required table exists
    expected_tables = {"devices", "labels", "topology_snapshots", "config"}
    cur = conn.execute("SELECT name FROM sqlite_master WHERE type='table';")
    tables = {row[0] for row in cur.fetchall()}
    assert expected_tables.issubset(tables)

    # Restore env
    if original_env is None:
        del db.os.environ["SQLITE_PATH"]
    else:
        db.os.environ["SQLITE_PATH"] = original_env
