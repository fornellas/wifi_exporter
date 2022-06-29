package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"gopkg.in/alecthomas/kingpin.v2"

	"github.com/fornellas/wifi_exporter/internal/network_manager"
	"github.com/fornellas/wifi_exporter/internal/wpa_supplicant"
	"github.com/fornellas/wifi_exporter/wifi"
)

var (
	address           = kingpin.Flag("server", "server address").Default(":8034").String()
	wifiBackend       = kingpin.Flag("wifi.backend", "Either wpa_supplicant or NetworkManager").Default("NetworkManager").String()
	wifiScanTimeoutMs = kingpin.Flag("wifi.scan.timeout_ms", "timeout scanning for wifi networks in milliseconds").Default("10000").Int()
)

func escape(s string) string {
	return strings.ReplaceAll(s, `"`, `\"`)
}

func getWifiClient() (wifi.Client, error) {
	switch *wifiBackend {
	case "NetworkManager":
		nm, err := network_manager.NewNetworkManager()
		if err != nil {
			return nil, err
		}
		return nm, nil
	case "wpa_supplicant":
		return wpa_supplicant.NewWPASupplicant(), nil
	default:
		return nil, fmt.Errorf("Invalid wifi backend: %s", *wifiBackend)
	}
}

func getwifiScanTimeoutDuration() time.Duration {
	return time.Duration(*wifiScanTimeoutMs) * time.Millisecond
}

func metricsHandler(w http.ResponseWriter, req *http.Request) {
	log.Printf("< %s GET /metrics", req.RemoteAddr)

	wifiClient, err := getWifiClient()
	if err != nil {
		log.Printf("> %s GET /metrics 500: %s", req.RemoteAddr, err.Error())
		w.WriteHeader(500)
		fmt.Fprintf(w, "Failed to get wifi client")
		return
	}

	ifaces, err := wifiClient.GetInterfaces()
	if err != nil {
		log.Printf("> %s GET /metrics 500: %s", req.RemoteAddr, err.Error())
		w.WriteHeader(500)
		fmt.Fprintf(w, "Failed to list interfaces")
		return
	}

	scanResultsMap := make(map[wifi.Iface][]wifi.ScanResult)
	scanStartTime := time.Now()
	for _, iface := range ifaces {
		scanResults, err := wifiClient.Scan(iface, getwifiScanTimeoutDuration())
		if err != nil {
			log.Printf("> %s GET /metrics 500: %s", req.RemoteAddr, err.Error())
			w.WriteHeader(500)
			fmt.Fprintf(w, "Failed to scan")
			return
		}
		scanResultsMap[iface] = scanResults
	}
	scanDuration := time.Since(scanStartTime)

	log.Printf("> %s GET /metrics 200", req.RemoteAddr)
	fmt.Fprintf(w, "wifi_scan_duration_seconds %f\n", float64(scanDuration.Milliseconds())/1000.0)
	for iface, scanResultList := range scanResultsMap {
		for _, scanResult := range scanResultList {
			fmt.Fprintf(
				w,
				"wifi_signal_db{interface=\"%s\",BSSID=\"%s\",SSID=\"%s\",frequency_MHz=\"%d\",frequency_band=\"%s\",channel=\"%s\"} %d\n",
				escape(iface.Name()),
				escape(scanResult.BSSID.String()),
				escape(scanResult.SSID),
				scanResult.FrequencyMHz,
				escape(scanResult.FrequencyBand()),
				escape(fmt.Sprintf("%d", scanResult.Channel())),
				scanResult.SignalStrengthdBm,
			)
			fmt.Fprintf(
				w,
				"wifi_channel{interface=\"%s\",BSSID=\"%s\",SSID=\"%s\",frequency_band=\"%s\"} %d\n",
				escape(iface.Name()),
				escape(scanResult.BSSID.String()),
				escape(scanResult.SSID),
				escape(scanResult.FrequencyBand()),
				scanResult.Channel(),
			)
		}
	}
}

func main() {

	kingpin.Parse()

	http.HandleFunc("/metrics", metricsHandler)

	log.Printf("Listening at %s", *address)
	if err := http.ListenAndServe(*address, nil); err != http.ErrServerClosed {
		log.Fatalf("Failed to start server: %s", err.Error())
	}
}
