
# Run a docker image



```bash
$ ./check_prometheus mode ping --address http://localhost:9090
OK - Version: 3.8.0, Instance localhost:9090|'duration'=0.001691195s;;;0;
```

# Target health testing

```bash
$ ./check_prometheus mode targets_health --address http://localhost:9090 -w 0.9: -c 0.5:
OK - There are 1 healthy and 0 unhealthy targets|'localhost:9090'=0;;;; 'health_rate'=1;0.9:;0.5:;;1 'targets'=1;;;0;
```

# Query Encoding testing

```bash
$ ./check_prometheus -t 9999 mode query -q 'up' --address http://localhost:9090 --insecure
OK - Query: 'up'|'{__name__="up", app="prometheus", instance="localhost:9090", job="prometheus"}'=1;;;;
```

```bash
$ ./check_prometheus -t 9999 mode query -q 'up' --query-encoding url --address http://localhost:9090 --insecure
OK - Query: 'up'|'{__name__="up", app="prometheus", instance="localhost:9090", job="prometheus"}'=1;;;;
```


```bash
$ ./check_prometheus -t 9999 mode query -q 'dXA=' --query-encoding base64 --address http://localhost:9090 --insecure
OK - Query: 'up'|'{__name__="up", app="prometheus", instance="localhost:9090", job="prometheus"}'=1;;;;
```


# Query Counter


```bash
$ ./check_prometheus mode query --address http://localhost:9090 -q 'rate(prometheus_http_requests_total[5m])'
OK - Query: 'rate(prometheus_http_requests_total[5m])'|'{app="prometheus", code="200", handler="/", instance="localhost:9090", job="prometheus"}'=0;;;; '{app="prometheus", code="200", handler="/-/healthy", instance="localhost:9090", job="prometheus"}'=0;;;; '{app="prometheus", code="200", handler="/-/quit", instance="localhost:9090", job="prometheus"}'=0;;;; '{app="prometheus", code="200", handler="/-/ready", instance="localhost:9090", job="prometheus"}'=0;;;; '{app="prometheus", code="200", handler="/-/reload", instance="localhost:9090", job="prometheus"}'=0;;;; '{app="prometheus", code="200", handler="/alertmanager-discovery", instance="localhost:9090", job="prometheus"}'=0;;;; '{app="prometheus", code="200", handler="/alerts", instance="localhost:9090", job="prometheus"}'=0;;;; '{app="prometheus", code="200", handler="/api/v1/*path", instance="localhost:9090", job="prometheus"}'=0;;;; '{app="prometheus", code="200", handler="/api/v1/admin/tsdb/clean_tombstones", instance="localhost:9090", job="prometheus"}'=0;;;; '{app="prometheus", code="200", handler="/api/v1/admin/tsdb/delete_series", instance="localhost:9090", job="prometheus"}'=0;;;; '{app="prometheus", code="200", handler="/api/v1/admin/tsdb/snapshot", instance="localhost:9090", job="prometheus"}'=0;;;; '{app="prometheus", code="200", handler="/api/v1/alertmanagers", instance="localhost:9090", job="prometheus"}'=0;;;; '{app="prometheus", code="200", handler="/api/v1/alerts", instance="localhost:9090", job="prometheus"}'=0;;;; '{app="prometheus", code="200", handler="/api/v1/features", instance="localhost:9090", job="prometheus"}'=0;;;; '{app="prometheus", code="200", handler="/api/v1/format_query", instance="localhost:9090", job="prometheus"}'=0;;;; '{app="prometheus", code="200", handler="/api/v1/label/:name/values", instance="localhost:9090", job="prometheus"}'=0;;;; '{app="prometheus", code="200", handler="/api/v1/labels", instance="localhost:9090", job="prometheus"}'=0;;;; '{app="prometheus", code="200", handler="/api/v1/metadata", instance="localhost:9090", job="prometheus"}'=0;;;; '{app="prometheus", code="200", handler="/api/v1/notifications", instance="localhost:9090", job="prometheus"}'=0;;;; '{app="prometheus", code="200", handler="/api/v1/notifications/live", instance="localhost:9090", job="prometheus"}'=0;;;; '{app="prometheus", code="200", handler="/api/v1/otlp/v1/metrics", instance="localhost:9090", job="prometheus"}'=0;;;; '{app="prometheus", code="200", handler="/api/v1/parse_query", instance="localhost:9090", job="prometheus"}'=0;;;; '{app="prometheus", code="200", handler="/api/v1/query", instance="localhost:9090", job="prometheus"}'=0.021052631578947364;;;; '{app="prometheus", code="200", handler="/api/v1/query_exemplars", instance="localhost:9090", job="prometheus"}'=0;;;; '{app="prometheus", code="200", handler="/api/v1/query_range", instance="localhost:9090", job="prometheus"}'=0;;;; '{app="prometheus", code="200", handler="/api/v1/read", instance="localhost:9090", job="prometheus"}'=0;;;; '{app="prometheus", code="200", handler="/api/v1/rules", instance="localhost:9090", job="prometheus"}'=0;;;; '{app="prometheus", code="200", handler="/api/v1/scrape_pools", instance="localhost:9090", job="prometheus"}'=0;;;; '{app="prometheus", code="200", handler="/api/v1/series", instance="localhost:9090", job="prometheus"}'=0;;;; '{app="prometheus", code="200", handler="/api/v1/status/buildinfo", instance="localhost:9090", job="prometheus"}'=0;;;; '{app="prometheus", code="200", handler="/api/v1/status/config", instance="localhost:9090", job="prometheus"}'=0;;;; '{app="prometheus", code="200", handler="/api/v1/status/flags", instance="localhost:9090", job="prometheus"}'=0;;;; '{app="prometheus", code="200", handler="/api/v1/status/runtimeinfo", instance="localhost:9090", job="prometheus"}'=0;;;; '{app="prometheus", code="200", handler="/api/v1/status/tsdb", instance="localhost:9090", job="prometheus"}'=0;;;; '{app="prometheus", code="200", handler="/api/v1/status/tsdb/blocks", instance="localhost:9090", job="prometheus"}'=0;;;; '{app="prometheus", code="200", handler="/api/v1/status/walreplay", instance="localhost:9090", job="prometheus"}'=0;;;; '{app="prometheus", code="200", handler="/api/v1/targets", instance="localhost:9090", job="prometheus"}'=0.003508771929824561;;;; '{app="prometheus", code="200", handler="/api/v1/targets/metadata", instance="localhost:9090", job="prometheus"}'=0;;;; '{app="prometheus", code="200", handler="/api/v1/targets/relabel_steps", instance="localhost:9090", job="prometheus"}'=0;;;; '{app="prometheus", code="200", handler="/api/v1/write", instance="localhost:9090", job="prometheus"}'=0;;;; '{app="prometheus", code="200", handler="/assets/*filepath", instance="localhost:9090", job="prometheus"}'=0;;;; '{app="prometheus", code="200", handler="/classic/static/*filepath", instance="localhost:9090", job="prometheus"}'=0;;;; '{app="prometheus", code="200", handler="/config", instance="localhost:9090", job="prometheus"}'=0;;;; '{app="prometheus", code="200", handler="/consoles/*filepath", instance="localhost:9090", job="prometheus"}'=0;;;; '{app="prometheus", code="200", handler="/debug/*subpath", instance="localhost:9090", job="prometheus"}'=0;;;; '{app="prometheus", code="200", handler="/favicon.ico", instance="localhost:9090", job="prometheus"}'=0;;;; '{app="prometheus", code="200", handler="/favicon.svg", instance="localhost:9090", job="prometheus"}'=0;;;; '{app="prometheus", code="200", handler="/federate", instance="localhost:9090", job="prometheus"}'=0;;;; '{app="prometheus", code="200", handler="/flags", instance="localhost:9090", job="prometheus"}'=0;;;; '{app="prometheus", code="200", handler="/graph", instance="localhost:9090", job="prometheus"}'=0;;;; '{app="prometheus", code="200", handler="/manifest.json", instance="localhost:9090", job="prometheus"}'=0;;;; '{app="prometheus", code="200", handler="/metrics", instance="localhost:9090", job="prometheus"}'=0.06666666666666667;;;; '{app="prometheus", code="200", handler="/query", instance="localhost:9090", job="prometheus"}'=0;;;; '{app="prometheus", code="200", handler="/rules", instance="localhost:9090", job="prometheus"}'=0;;;; '{app="prometheus", code="200", handler="/service-discovery", instance="localhost:9090", job="prometheus"}'=0;;;; '{app="prometheus", code="200", handler="/status", instance="localhost:9090", job="prometheus"}'=0;;;; '{app="prometheus", code="200", handler="/targets", instance="localhost:9090", job="prometheus"}'=0;;;; '{app="prometheus", code="200", handler="/tsdb-status", instance="localhost:9090", job="prometheus"}'=0;;;; '{app="prometheus", code="200", handler="/version", instance="localhost:9090", job="prometheus"}'=0;;;; '{app="prometheus", code="400", handler="/api/v1/query", instance="localhost:9090", job="prometheus"}'=0;;;;
```


# Query Gauge

```bash
$ ./check_prometheus mode query --address http://localhost:9090 -q 'prometheus_tsdb_head_series'
OK - Query: 'prometheus_tsdb_head_series'|'{__name__="prometheus_tsdb_head_series", app="prometheus", instance="localhost:9090", job="prometheus"}'=928;;;;
```


# Query Histograms


```bash
./check_prometheus mode query --address http://localhost:9090 -q 'histogram_quantile(0.99, sum by (le) (rate(prometheus_http_request_duration_seconds_bucket[5m])))'
```




# Aliasing


```bash
./check_prometheus mode query --address http://localhost:9090 -q 'prometheus_http_requests_total{handler="/metrics"}' -a 'The handler {{.handler}} has received {{.xvalue}} requests.'
```
