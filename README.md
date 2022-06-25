# wifi_exporter

Prometheus exporter that exposes Wireless interfaces information.

## Metrics

| Name  | Description | Labels |
| -- | -- | -- |
| wifi_signal_db | BSS signal in dB (RSSI) | `interface`, `BSSID`, `SSID`, `frequency_MHz`, `frequency_band`, `channel`, `flags` |