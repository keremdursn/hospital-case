#!/bin/bash
# Bu script, Prometheus metriklerini kullanan örnek bir Grafana dashboard'unu otomatik olarak import eder.
# GRAFANA_URL ve API_KEY değerlerini kendi ortamınıza göre değiştirin.

GRAFANA_URL="http://localhost:3000"
API_KEY="YOUR_GRAFANA_API_KEY"

# Örnek dashboard JSON'u (örnek olarak bir panel ile)
cat <<EOF > prometheus_dashboard.json
{
  "dashboard": {
    "id": null,
    "uid": null,
    "title": "Go Hospital Case Metrics",
    "tags": [ "prometheus", "go" ],
    "timezone": "browser",
    "schemaVersion": 16,
    "version": 0,
    "panels": [
      {
        "type": "graph",
        "title": "HTTP Requests",
        "targets": [
          {
            "expr": "http_requests_total",
            "legendFormat": "{{method}} {{path}} {{status}}",
            "refId": "A"
          }
        ],
        "datasource": "Prometheus",
        "id": 1
      }
    ]
  },
  "overwrite": true
}
EOF

# Dashboard'u import et
curl -X POST \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $API_KEY" \
  -d @prometheus_dashboard.json \
  $GRAFANA_URL/api/dashboards/db

echo "Dashboard import edildi. Grafana arayüzünden kontrol edebilirsiniz." 