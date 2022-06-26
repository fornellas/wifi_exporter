package wireless

import (
	"errors"
	"fmt"
	"log"
	"net"
	"os"
	"path"
	"path/filepath"
	"pifke.org/wpasupplicant"
	"reflect"
	"sync"
	"time"
)

var scanMutex sync.Mutex

func GetWirelessInterfaceNames() (ifNames []string, err error) {
	log.Printf("Listing wireless interfaces")
	matches, err := filepath.Glob("/sys/class/net/*")
	if err != nil {
		return nil, err
	}

	for _, ifName := range matches {
		_, err := os.Stat(path.Join(ifName, "wireless"))
		if err != nil {
			if errors.Is(err, os.ErrNotExist) {
				continue
			}
			return nil, err
		}
		ifNames = append(ifNames, path.Base(ifName))
	}

	return ifNames, nil
}

type ScanResult struct {
	// IfName is the interface name
	IfName string
	// BSSID is the MAC address of the BSS.
	BSSID net.HardwareAddr
	// SSID is the SSID of the BSS.
	SSID string
	// Frequency is the frequency, in Mhz, of the BSS.
	Frequency int
	// RSSI is the received signal strength, in dB, of the BSS.
	RSSI int
	// Flags is an array of flags, in string format, returned by the
	// wpa_supplicant SCAN_RESULTS command.  Future versions of this code
	// will parse these into something more meaningful.
	Flags []string
}

func (s *ScanResult) FrequencyBand() string {
	if s.Frequency >= 2412 && s.Frequency <= 2484 {
		return "2.4GHz"
	}
	if s.Frequency >= 5035 && s.Frequency <= 5980 {
		return "5GHz"
	}
	return "Unknown"
}

func (s *ScanResult) Channel() int {
	// 2.4GHz https://en.wikipedia.org/wiki/List_of_WLAN_channels#2.4_GHz_(802.11b/g/n/ax)
	if s.Frequency == 2412 {
		return 1
	}
	if s.Frequency == 2417 {
		return 2
	}
	if s.Frequency == 2422 {
		return 3
	}
	if s.Frequency == 2427 {
		return 4
	}
	if s.Frequency == 2432 {
		return 5
	}
	if s.Frequency == 2437 {
		return 6
	}
	if s.Frequency == 2442 {
		return 7
	}
	if s.Frequency == 2447 {
		return 8
	}
	if s.Frequency == 2452 {
		return 9
	}
	if s.Frequency == 2457 {
		return 10
	}
	if s.Frequency == 2462 {
		return 11
	}
	if s.Frequency == 2467 {
		return 12
	}
	if s.Frequency == 2472 {
		return 13
	}
	if s.Frequency == 2484 {
		return 14
	}

	// 5GHz https://en.wikipedia.org/wiki/List_of_WLAN_channels#5_GHz_(802.11a/h/j/n/ac/ax)
	if s.Frequency == 5035 {
		return 7
	}
	if s.Frequency == 5040 {
		return 8
	}
	if s.Frequency == 5045 {
		return 9
	}
	if s.Frequency == 5055 {
		return 11
	}
	if s.Frequency == 5060 {
		return 12
	}
	if s.Frequency == 5080 {
		return 16
	}
	if s.Frequency == 5160 {
		return 32
	}
	if s.Frequency == 5170 {
		return 34
	}
	if s.Frequency == 5180 {
		return 36
	}
	if s.Frequency == 5190 {
		return 38
	}
	if s.Frequency == 5200 {
		return 40
	}
	if s.Frequency == 5210 {
		return 42
	}
	if s.Frequency == 5220 {
		return 44
	}
	if s.Frequency == 5230 {
		return 46
	}
	if s.Frequency == 5240 {
		return 48
	}
	if s.Frequency == 5250 {
		return 50
	}
	if s.Frequency == 5260 {
		return 52
	}
	if s.Frequency == 5270 {
		return 54
	}
	if s.Frequency == 5280 {
		return 56
	}
	if s.Frequency == 5290 {
		return 58
	}
	if s.Frequency == 5300 {
		return 60
	}
	if s.Frequency == 5310 {
		return 62
	}
	if s.Frequency == 5320 {
		return 64
	}
	if s.Frequency == 5340 {
		return 68
	}
	if s.Frequency == 5480 {
		return 96
	}
	if s.Frequency == 5500 {
		return 100
	}
	if s.Frequency == 5510 {
		return 102
	}
	if s.Frequency == 5520 {
		return 104
	}
	if s.Frequency == 5530 {
		return 106
	}
	if s.Frequency == 5540 {
		return 108
	}
	if s.Frequency == 5550 {
		return 110
	}
	if s.Frequency == 5560 {
		return 112
	}
	if s.Frequency == 5570 {
		return 114
	}
	if s.Frequency == 5580 {
		return 116
	}
	if s.Frequency == 5590 {
		return 118
	}
	if s.Frequency == 5600 {
		return 120
	}
	if s.Frequency == 5610 {
		return 122
	}
	if s.Frequency == 5620 {
		return 124
	}
	if s.Frequency == 5630 {
		return 126
	}
	if s.Frequency == 5640 {
		return 128
	}
	if s.Frequency == 5660 {
		return 132
	}
	if s.Frequency == 5670 {
		return 134
	}
	if s.Frequency == 5680 {
		return 136
	}
	if s.Frequency == 5690 {
		return 138
	}
	if s.Frequency == 5700 {
		return 140
	}
	if s.Frequency == 5710 {
		return 142
	}
	if s.Frequency == 5720 {
		return 144
	}
	if s.Frequency == 5745 {
		return 149
	}
	if s.Frequency == 5755 {
		return 151
	}
	if s.Frequency == 5765 {
		return 153
	}
	if s.Frequency == 5775 {
		return 155
	}
	if s.Frequency == 5785 {
		return 157
	}
	if s.Frequency == 5795 {
		return 159
	}
	if s.Frequency == 5805 {
		return 161
	}
	if s.Frequency == 5815 {
		return 163
	}
	if s.Frequency == 5825 {
		return 165
	}
	if s.Frequency == 5835 {
		return 167
	}
	if s.Frequency == 5845 {
		return 169
	}
	if s.Frequency == 5855 {
		return 171
	}
	if s.Frequency == 5865 {
		return 173
	}
	if s.Frequency == 5875 {
		return 175
	}
	if s.Frequency == 5885 {
		return 177
	}
	if s.Frequency == 5910 {
		return 182
	}
	if s.Frequency == 5915 {
		return 183
	}
	if s.Frequency == 5920 {
		return 184
	}
	if s.Frequency == 5935 {
		return 187
	}
	if s.Frequency == 5940 {
		return 188
	}
	if s.Frequency == 5945 {
		return 189
	}
	if s.Frequency == 5960 {
		return 192
	}
	if s.Frequency == 5980 {
		return 196
	}

	return 0
}

func Scan(timeout time.Duration) ([]ScanResult, error) {
	ifNames, err := GetWirelessInterfaceNames()
	if err != nil {
		return nil, fmt.Errorf("Failed to list wireless interfaces: %s", err.Error())
	}

	scanMutex.Lock()
	defer scanMutex.Unlock()

	scanResults := []ScanResult{}
	for _, ifName := range ifNames {
		log.Printf("WPA Supplicant connection to %s", ifName)
		conn, err := wpasupplicant.Unixgram(ifName)
		if err != nil {
			return nil, fmt.Errorf("Failed to connect to WPA Supplicant %s: %s", ifName, err.Error())
		}
		defer conn.Close()
		// ScanResults() may hang forever without this
		conn.SetTimeout(timeout)

		log.Printf("Scanning %s", ifName)
		timeoutCh := time.After(timeout)
		err = conn.Scan()
		if err != nil {
			return nil, fmt.Errorf("Failed to scan %s: %s (%v)", ifName, err.Error(), reflect.TypeOf(err))
		}

		scanComplete := false
		for {
			select {
			case wpaEvent := <-conn.EventQueue():
				if wpaEvent.Event == "SCAN-RESULTS" {
					scanComplete = true
					break
				}
			case <-timeoutCh:
				return nil, fmt.Errorf("Scan %s timeout after %s", ifName, timeout)
			}
			if scanComplete {
				break
			}
		}

		connScanResults, errs := conn.ScanResults()
		if errs != nil {
			return nil, fmt.Errorf("Failed to get scan results %s: %s", ifName, errs)
		}

		for _, connScanResult := range connScanResults {
			scanResults = append(
				scanResults,
				ScanResult{
					IfName:    ifName,
					BSSID:     connScanResult.BSSID(),
					SSID:      connScanResult.SSID(),
					Frequency: connScanResult.Frequency(),
					RSSI:      connScanResult.RSSI(),
					Flags:     connScanResult.Flags(),
				},
			)
		}
	}

	return scanResults, nil
}
