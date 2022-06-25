# wifi_exporter

Prometheus exporter that exposes Wireless interfaces information.

## Metrics

| Name  | Description | Labels |
| -- | -- | -- |
| wifi_scan | Scan result of available networks | `interface`, `BSSID`, `SSID`, `frequency_MHz`, `RSSI`, `flags` |