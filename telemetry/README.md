# Telemetry

To run the telemetry infrastructure you can run `docker-compose up` in this directory.
This will setup and run the following components.

## Components

### Prometheus

Prometheus handles all the metrics
Currently configured to simply scrape pushgateway.

Point grafana to <http://localhost:9090>

UI <http://localhost:9090>

### Push gateway

Push gateway provides a gateway for prometheus metrics.

UI <http://localhost:9091>

### Loki

Loki handles storing logs.

Point grafana to <http://localhost:3100>

### Jaeger

Jaeger handles storing traces.

API to push traces: http://localhost:4317

Point grafana to <http://localhost:16686>

UI <http://localhost:16686>
