// Task: downsample bandwidth measurement older than 7d into 15m aggregates
import "influxdata/influxdb/tasks"

option task = {
  name: "downsample_bandwidth_15m",
  every: 24h,
  offset: 10m,
  description: "Aggregate 1m bandwidth points to 15m buckets for data older than 7 days",
}

from(bucket: "network_metrics")
  |> range(start: -90d, stop: -7d)
  |> filter(fn: (r) => r._measurement == "bandwidth")
  |> aggregateWindow(every: 15m, fn: max, createEmpty: false)
  |> to(bucket: "network_metrics", org: "home-net", fieldFn: (r) => ({rx_bytes: r.rx_bytes, tx_bytes: r.tx_bytes}), tagColumns: ["mac", "interface"])
