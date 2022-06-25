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

func getFrequencyBand(frequency int) string {
	if frequency >= 2412 && frequency <= 2484 {
		return "2.4GHz"
	}
	if frequency >= 5035 && frequency <= 5980 {
		return "5GHz"
	}
	return "Unknown"
}

func getChannel(frequency int) int {
	// 2.4GHz https://en.wikipedia.org/wiki/List_of_WLAN_channels#2.4_GHz_(802.11b/g/n/ax)
	if frequency == 2412 {
		return 1
	}
	if frequency == 2417 {
		return 2
	}
	if frequency == 2422 {
		return 3
	}
	if frequency == 2427 {
		return 4
	}
	if frequency == 2432 {
		return 5
	}
	if frequency == 2437 {
		return 6
	}
	if frequency == 2442 {
		return 7
	}
	if frequency == 2447 {
		return 8
	}
	if frequency == 2452 {
		return 9
	}
	if frequency == 2457 {
		return 10
	}
	if frequency == 2462 {
		return 11
	}
	if frequency == 2467 {
		return 12
	}
	if frequency == 2472 {
		return 13
	}
	if frequency == 2484 {
		return 14
	}

	// 5GHz https://en.wikipedia.org/wiki/List_of_WLAN_channels#5_GHz_(802.11a/h/j/n/ac/ax)
	if frequency == 5035 {
		return 7
	}
	if frequency == 5040 {
		return 8
	}
	if frequency == 5045 {
		return 9
	}
	if frequency == 5055 {
		return 11
	}
	if frequency == 5060 {
		return 12
	}
	if frequency == 5080 {
		return 16
	}
	if frequency == 5160 {
		return 32
	}
	if frequency == 5170 {
		return 34
	}
	if frequency == 5180 {
		return 36
	}
	if frequency == 5190 {
		return 38
	}
	if frequency == 5200 {
		return 40
	}
	if frequency == 5210 {
		return 42
	}
	if frequency == 5220 {
		return 44
	}
	if frequency == 5230 {
		return 46
	}
	if frequency == 5240 {
		return 48
	}
	if frequency == 5250 {
		return 50
	}
	if frequency == 5260 {
		return 52
	}
	if frequency == 5270 {
		return 54
	}
	if frequency == 5280 {
		return 56
	}
	if frequency == 5290 {
		return 58
	}
	if frequency == 5300 {
		return 60
	}
	if frequency == 5310 {
		return 62
	}
	if frequency == 5320 {
		return 64
	}
	if frequency == 5340 {
		return 68
	}
	if frequency == 5480 {
		return 96
	}
	if frequency == 5500 {
		return 100
	}
	if frequency == 5510 {
		return 102
	}
	if frequency == 5520 {
		return 104
	}
	if frequency == 5530 {
		return 106
	}
	if frequency == 5540 {
		return 108
	}
	if frequency == 5550 {
		return 110
	}
	if frequency == 5560 {
		return 112
	}
	if frequency == 5570 {
		return 114
	}
	if frequency == 5580 {
		return 116
	}
	if frequency == 5590 {
		return 118
	}
	if frequency == 5600 {
		return 120
	}
	if frequency == 5610 {
		return 122
	}
	if frequency == 5620 {
		return 124
	}
	if frequency == 5630 {
		return 126
	}
	if frequency == 5640 {
		return 128
	}
	if frequency == 5660 {
		return 132
	}
	if frequency == 5670 {
		return 134
	}
	if frequency == 5680 {
		return 136
	}
	if frequency == 5690 {
		return 138
	}
	if frequency == 5700 {
		return 140
	}
	if frequency == 5710 {
		return 142
	}
	if frequency == 5720 {
		return 144
	}
	if frequency == 5745 {
		return 149
	}
	if frequency == 5755 {
		return 151
	}
	if frequency == 5765 {
		return 153
	}
	if frequency == 5775 {
		return 155
	}
	if frequency == 5785 {
		return 157
	}
	if frequency == 5795 {
		return 159
	}
	if frequency == 5805 {
		return 161
	}
	if frequency == 5815 {
		return 163
	}
	if frequency == 5825 {
		return 165
	}
	if frequency == 5835 {
		return 167
	}
	if frequency == 5845 {
		return 169
	}
	if frequency == 5855 {
		return 171
	}
	if frequency == 5865 {
		return 173
	}
	if frequency == 5875 {
		return 175
	}
	if frequency == 5885 {
		return 177
	}
	if frequency == 5910 {
		return 182
	}
	if frequency == 5915 {
		return 183
	}
	if frequency == 5920 {
		return 184
	}
	if frequency == 5935 {
		return 187
	}
	if frequency == 5940 {
		return 188
	}
	if frequency == 5945 {
		return 189
	}
	if frequency == 5960 {
		return 192
	}
	if frequency == 5980 {
		return 196
	}

	return 0
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
			"wifi_signal_db{interface=\"%s\",BSSID=\"%s\",SSID=\"%s\",frequency_MHz=\"%s\",frequency_band=\"%s\",channel=\"%s\",flags=\"%s\"} %d\n",
			escape(scanResult.IfName),
			escape(scanResult.BSSID.String()),
			escape(scanResult.SSID),
			escape(strconv.Itoa(scanResult.Frequency)),
			escape(getFrequencyBand(scanResult.Frequency)),
			escape(fmt.Sprintf("%d", getChannel(scanResult.Frequency))),
			escape(fmt.Sprintf("%s", scanResult.Flags)),
			scanResult.RSSI,
		)
		fmt.Fprintf(
			w,
			"wifi_channel{interface=\"%s\",BSSID=\"%s\",SSID=\"%s\",frequency_band=\"%s\",flags=\"%s\"} %d\n",
			escape(scanResult.IfName),
			escape(scanResult.BSSID.String()),
			escape(scanResult.SSID),
			escape(getFrequencyBand(scanResult.Frequency)),
			escape(fmt.Sprintf("%s", scanResult.Flags)),
			getChannel(scanResult.Frequency),
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
