package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"gopkg.in/alecthomas/kingpin.v2"

	"github.com/fornellas/wifi_exporter/internal/wpa_supplicant"
	"github.com/fornellas/wifi_exporter/wifi"
)

var (
	address              = kingpin.Flag("server", "server address").Default(":8034").String()
	wpaSupplicantTimeout = kingpin.Flag("wifi.wpa_supplicant.timeout_ms", "timeout when talking to WPA Supplicant in milliseconds").Default("10000").Int()
)

func escape(s string) string {
	return strings.ReplaceAll(s, `"`, `\"`)
}

func metricsHandler(w http.ResponseWriter, req *http.Request) {
	log.Printf("< %s GET /metrics", req.RemoteAddr)

	wifiClient := wpa_supplicant.NewWPASupplicant()

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
		scanResults, err := wifiClient.Scan(iface, time.Duration(*wpaSupplicantTimeout)*time.Millisecond)
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
				"wifi_signal_db{interface=\"%s\",BSSID=\"%s\",SSID=\"%s\",frequency_MHz=\"%s\",frequency_band=\"%s\",channel=\"%s\",flags=\"%s\"} %d\n",
				escape(iface.Name()),
				escape(scanResult.BSSID.String()),
				escape(scanResult.SSID),
				escape(strconv.Itoa(scanResult.Frequency)),
				escape(scanResult.FrequencyBand()),
				escape(fmt.Sprintf("%d", scanResult.Channel())),
				escape(fmt.Sprintf("%s", scanResult.Flags)),
				scanResult.RSSI,
			)
			fmt.Fprintf(
				w,
				"wifi_channel{interface=\"%s\",BSSID=\"%s\",SSID=\"%s\",frequency_band=\"%s\",flags=\"%s\"} %d\n",
				escape(iface.Name()),
				escape(scanResult.BSSID.String()),
				escape(scanResult.SSID),
				escape(scanResult.FrequencyBand()),
				escape(fmt.Sprintf("%s", scanResult.Flags)),
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
