package main

import (
	"fmt"
	"github.com/fornellas/wifi_exporter/wireless"
	"gopkg.in/alecthomas/kingpin.v2"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

var (
	address               = kingpin.Flag("server", "server address").Default(":8034").String()
	wpaSupplicantTimeout = kingpin.Flag("wireless.wpa_supplicant.timeout_ms", "timeout when talking to WPA Supplicant in milliseconds").Default("5000").Int()
)

func escape(s string) string {
	return strings.ReplaceAll(s, `"`, `\"`)
}

func metricsHandler(w http.ResponseWriter, req *http.Request) {
	log.Printf("< %s GET /metrics", req.RemoteAddr)
	scanStartTime := time.Now()
	scanResults, err := wireless.Scan(time.Duration(*wpaSupplicantTimeout) * time.Millisecond)
	scanDuration := time.Now().Sub(scanStartTime)
	if err != nil {
		log.Printf("> %s GET /metrics 500: %s", req.RemoteAddr, err.Error())
		w.WriteHeader(500)
		fmt.Fprintf(w, "Failed to scan")
		return
	}
	log.Printf("> %s GET /metrics 200", req.RemoteAddr)
	fmt.Fprintf(w, "wifi_scan_duration_seconds %f\n", float64(scanDuration.Milliseconds())/1000.0)
	for _, scanResult := range scanResults {
		fmt.Fprintf(
			w,
			"wifi_signal_db{interface=\"%s\",BSSID=\"%s\",SSID=\"%s\",frequency_MHz=\"%s\",frequency_band=\"%s\",channel=\"%s\",flags=\"%s\"} %d\n",
			escape(scanResult.IfName),
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
			escape(scanResult.IfName),
			escape(scanResult.BSSID.String()),
			escape(scanResult.SSID),
			escape(scanResult.FrequencyBand()),
			escape(fmt.Sprintf("%s", scanResult.Flags)),
			scanResult.Channel(),
		)
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
