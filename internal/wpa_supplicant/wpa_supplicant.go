package wpa_supplicant

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path"
	"path/filepath"
	"reflect"
	"sync"
	"time"

	"pifke.org/wpasupplicant"

	"github.com/fornellas/wifi_exporter/wifi"
)

type WPASupplicant struct {
	wifi.Client
}

func NewWPASupplicant() WPASupplicant {
	return WPASupplicant{}
}

var scanMutex sync.Mutex

func (_ *WPASupplicant) GetInterfaces() (ifaces []wifi.Iface, err error) {
	log.Printf("Listing wireless interfaces")
	matches, err := filepath.Glob("/sys/class/net/*")
	if err != nil {
		return nil, err
	}

	for _, iface := range matches {
		_, err := os.Stat(path.Join(iface, "wireless"))
		if err != nil {
			if errors.Is(err, os.ErrNotExist) {
				continue
			}
			return nil, err
		}
		ifaces = append(ifaces, wifi.NewIface(path.Base(iface)))
	}

	return ifaces, nil
}

func (ws *WPASupplicant) Scan(iface wifi.Iface, timeout time.Duration) ([]wifi.ScanResult, error) {
	// TODO lock by interface
	scanMutex.Lock()
	defer scanMutex.Unlock()

	scanResults := []wifi.ScanResult{}
	log.Printf("WPA Supplicant connection to %s", iface)
	conn, err := wpasupplicant.Unixgram(iface.Name())
	if err != nil {
		return nil, fmt.Errorf("Failed to connect to WPA Supplicant %s: %s", iface, err.Error())
	}
	defer conn.Close()
	// ScanResults() may hang forever without this
	conn.SetTimeout(timeout)

	log.Printf("Scanning %s", iface)
	timeoutCh := time.After(timeout)
	err = conn.Scan()
	if err != nil {
		return nil, fmt.Errorf("Failed to scan %s: %s (%v)", iface, err.Error(), reflect.TypeOf(err))
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
			return nil, fmt.Errorf("Scan %s timeout after %s", iface, timeout)
		}
		if scanComplete {
			break
		}
	}

	connScanResults, errs := conn.ScanResults()
	if errs != nil {
		return nil, fmt.Errorf("Failed to get scan results %s: %s", iface, errs)
	}

	for _, connScanResult := range connScanResults {
		scanResults = append(
			scanResults,
			wifi.ScanResult{
				Iface:     iface,
				BSSID:     connScanResult.BSSID(),
				SSID:      connScanResult.SSID(),
				Frequency: connScanResult.Frequency(),
				RSSI:      connScanResult.RSSI(),
				Flags:     connScanResult.Flags(),
			},
		)
	}

	return scanResults, nil
}
