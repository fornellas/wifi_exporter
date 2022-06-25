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
	wirelessScanTimeoutMs = kingpin.Flag("wireless.scan.timeout_ms", "wireless scan timeout in milliseconds").Default("10000").Int()
)

func escape(s string) string {
	return strings.ReplaceAll(s, `"`, `\"`)
}

func metricsHandler(w http.ResponseWriter, req *http.Request) {
	log.Printf("< GET /metrics")
	scanResults, err := wireless.Scan(time.Duration(*wirelessScanTimeoutMs) * time.Millisecond)
	if err != nil {
		log.Printf("> GET /metrics 500")
		w.WriteHeader(500)
		fmt.Fprintf(w, "Failed to scan: %s", err.Error())
	}
	log.Printf("> GET /metrics 200")
	for _, scanResult := range scanResults {
		fmt.Fprintf(
			w,
			"wifi_scan{interface=\"%s\",BSSID=\"%s\",SSID=\"%s\",frequency_MHz=\"%s\",RSSI=\"%s\",flags=\"%s\"} 1\n",
			escape(scanResult.IfName),
			escape(scanResult.BSSID.String()),
			escape(scanResult.SSID),
			escape(strconv.Itoa(scanResult.Frequency)),
			escape(strconv.Itoa(scanResult.RSSI)),
			escape(fmt.Sprintf("%s", scanResult.Flags)),
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
