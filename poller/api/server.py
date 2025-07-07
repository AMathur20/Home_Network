from fastapi import FastAPI, HTTPException
from pydantic import BaseModel
import threading
from typing import List

from db import get_connection
from polling.unifi import poll_unifi
from polling.mikrotik import poll_mikrotik
from topology.engine import update_topology
from processor import process_data

app = FastAPI()

class Label(BaseModel):
    mac: str
    label: str

@app.post("/refresh-now")
def manual_refresh():
    def refresh():
        unifi_data = poll_unifi()
        mikrotik_data = poll_mikrotik()
        combined = unifi_data + mikrotik_data
        process_data(combined)
        update_topology(unifi_data, mikrotik_data)
    threading.Thread(target=refresh).start()
    return {"status": "refresh triggered"}

@app.get("/labels", response_model=List[Label])
def list_labels():
    conn = get_connection()
    cur = conn.execute("SELECT mac, label FROM labels ORDER BY label;")
    return [Label(**row) for row in cur.fetchall()]


@app.post("/labels", response_model=Label)
def upsert_label(label: Label):
    conn = get_connection()
    conn.execute(
        "INSERT INTO labels (mac, label) VALUES (?, ?) ON CONFLICT(mac) DO UPDATE SET label=excluded.label;",
        (label.mac.lower(), label.label),
    )
    conn.commit()
    return label


@app.delete("/labels/{mac}")
def delete_label(mac: str):
    conn = get_connection()
    cur = conn.execute("DELETE FROM labels WHERE mac = ?;", (mac.lower(),))
    conn.commit()
    if cur.rowcount == 0:
        raise HTTPException(status_code=404, detail="MAC not found")
    return {"deleted": mac}


from fastapi.responses import HTMLResponse

@app.get("/labels-ui", response_class=HTMLResponse)
def labels_ui():
    return """
<!doctype html>
<html lang='en'>
<head>
  <meta charset='utf-8'>
  <meta name='viewport' content='width=device-width, initial-scale=1'>
  <link href='https://cdn.jsdelivr.net/npm/bootstrap@5.3.3/dist/css/bootstrap.min.css' rel='stylesheet'>
  <title>Device Labels</title>
</head>
<body class='container py-4'>
  <h1 class='mb-4'>Device Labels</h1>
  <table class='table table-striped' id='labels-table'>
    <thead><tr><th>MAC</th><th>Label</th><th></th></tr></thead>
    <tbody></tbody>
  </table>
  <h3 class='mt-5'>Add / Update Label</h3>
  <form id='label-form' class='row g-3'>
    <div class='col-md-4'>
      <label class='form-label'>MAC</label>
      <input type='text' class='form-control' id='mac' required>
    </div>
    <div class='col-md-4'>
      <label class='form-label'>Label</label>
      <input type='text' class='form-control' id='label' required>
    </div>
    <div class='col-md-2 align-self-end'>
      <button type='submit' class='btn btn-primary'>Save</button>
    </div>
  </form>

<script>
async function loadLabels() {
  const resp = await fetch('/labels');
  const data = await resp.json();
  const tbody = document.querySelector('#labels-table tbody');
  tbody.innerHTML = '';
  data.forEach(row => {
    const tr = document.createElement('tr');
    tr.innerHTML = `<td>${row.mac}</td><td>${row.label}</td><td><button class='btn btn-sm btn-danger'>Delete</button></td>`;
    tr.querySelector('button').addEventListener('click', async () => {
      await fetch(`/labels/${row.mac}`, {method: 'DELETE'});
      loadLabels();
    });
    tbody.appendChild(tr);
  });
}

document.getElementById('label-form').addEventListener('submit', async (e) => {
  e.preventDefault();
  const mac = document.getElementById('mac').value;
  const label = document.getElementById('label').value;
  await fetch('/labels', {method: 'POST', headers: {'Content-Type': 'application/json'}, body: JSON.stringify({mac, label})});
  e.target.reset();
  loadLabels();
});

loadLabels();
</script>
</body>
</html>
"""


def start_api():
    import uvicorn
    uvicorn.run(app, host="0.0.0.0", port=5000)
