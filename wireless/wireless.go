package wireless

import (
	"errors"
	"fmt"
	"net"
	"os"
	"path"
	"path/filepath"
	"pifke.org/wpasupplicant"
	"reflect"
	"time"
	"sync"
	"log"
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
