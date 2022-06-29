package network_manager

import (
	"fmt"
	"log"
	"net"
	"time"

	"github.com/Wifx/gonetworkmanager"

	"github.com/fornellas/wifi_exporter/wifi"
)

type NetworkManager struct {
	wifi.Client
	gonm gonetworkmanager.NetworkManager
}

func NewNetworkManager() (NetworkManager, error) {
	var nm NetworkManager
	gonm, err := gonetworkmanager.NewNetworkManager()
	if err != nil {
		return nm, fmt.Errorf("Failed NewNetworkManager(): %s", err.Error())
	}
	nm.gonm = gonm
	return nm, nil
}

func (nm NetworkManager) GetInterfaces() (ifaces []wifi.Iface, err error) {
	log.Printf("Listing interfaces")
	devices, err := nm.gonm.GetPropertyAllDevices()
	if err != nil {
		return nil, fmt.Errorf("Failed GetPropertyAllDevices(): %s", err.Error())
	}
	for _, device := range devices {
		nmDeviceType, err := device.GetPropertyDeviceType()
		if err != nil {
			return nil, fmt.Errorf("Failed GetPropertyDeviceType(): %s", err.Error())
		}
		if nmDeviceType != gonetworkmanager.NmDeviceTypeWifi {
			continue
		}
		name, err := device.GetPropertyInterface()
		if err != nil {
			return nil, fmt.Errorf("Failed GetPropertyInterface(): %s", err.Error())
		}
		ifaces = append(ifaces, wifi.NewIface(name))
	}
	return ifaces, nil
}

func (nm NetworkManager) getDeviceByName(name string) (gonetworkmanager.Device, error) {
	devices, err := nm.gonm.GetPropertyAllDevices()
	if err != nil {
		return nil, fmt.Errorf("Failed GetPropertyAllDevices(): %s", err.Error())
	}
	for _, device := range devices {
		deviceName, err := device.GetPropertyInterface()
		if err != nil {
			return nil, fmt.Errorf("Failed GetPropertyInterface(): %s", err.Error())
		}
		if deviceName == name {
			return device, nil
		}
	}
	return nil, fmt.Errorf("Can not find device: %s", name)
}

// Extracted from
// https://github.com/NetworkManager/NetworkManager/blob/e6a33c04ebe1ac84e31628911e25bdfd7534dd3c/src/core/nm-core-utils.c#L4813-L4830
func getdBm(percent uint8) int {
	return int(((((-1 * (float64(percent) - 100.0)) * 60.0) / 100.0) / -1.0) - 40.0)
}

func (nm NetworkManager) Scan(iface wifi.Iface, timeout time.Duration) ([]wifi.ScanResult, error) {
	device, err := nm.getDeviceByName(iface.Name())
	if err != nil {
		return nil, fmt.Errorf("Failed getDeviceByName(): %s", err.Error())
	}
	deviceWireless, err := gonetworkmanager.NewDeviceWireless(device.GetPath())
	if err != nil {
		return nil, fmt.Errorf("Failed NewDeviceWireless(): %s", err.Error())
	}

	lastScan, err := deviceWireless.GetPropertyLastScan()
	if err != nil {
		return nil, fmt.Errorf("Failed GetPropertyLastScan(): %s", err.Error())
	}

	waitScanResult := true
	log.Printf("Requesting scan")
	err = deviceWireless.RequestScan()
	if err != nil {
		if err.Error() == "Scanning not allowed immediately following previous scan" {
			log.Printf("Recent scan available, moving on")
			waitScanResult = false
		} else if err.Error() == "Scanning not allowed while already scanning" {
			log.Printf("Scan already running")
		} else {
			return nil, fmt.Errorf("Failed RequestScan(): %v", err)
		}
	}

	if waitScanResult {
		log.Printf("Waiting for scan result")
		scanCompleteCh := make(chan error)
		go func() {
			timeoutCh := time.After(timeout)
			for {
				select {
				case <-time.After(100 * time.Millisecond):
					newLastScan, err := deviceWireless.GetPropertyLastScan()
					if err != nil {
						scanCompleteCh <- fmt.Errorf("Failed GetPropertyLastScan(): %s", err.Error())
						return
					}
					if newLastScan > lastScan {
						scanCompleteCh <- nil
						return
					}
				case <-timeoutCh:
					scanCompleteCh <- fmt.Errorf("Scan %s timeout after %s", iface, timeout)
					return
				}
			}
		}()

		timeoutCh := time.After(timeout)
		select {
		case err = <-scanCompleteCh:
			if err != nil {
				return nil, err
			}
			break
		case <-timeoutCh:
			return nil, fmt.Errorf("Scan %s timeout after %s", iface, timeout)
		}
	}

	log.Printf("Reading scan results")
	accessPoints, err := deviceWireless.GetAllAccessPoints()
	if err != nil {
		return nil, fmt.Errorf("Failed GetAllAccessPoints(): %s", err.Error())
	}
	scanResults := []wifi.ScanResult{}
	for _, accessPoint := range accessPoints {
		mac, err := accessPoint.GetPropertyHWAddress()
		if err != nil {
			return nil, fmt.Errorf("Failed GetPropertyHWAddress(): %s", err.Error())
		}
		bssid, err := net.ParseMAC(mac)
		if err != nil {
			return nil, fmt.Errorf("Failed ParseMAC(): %s", err.Error())
		}
		ssid, err := accessPoint.GetPropertySSID()
		if err != nil {
			return nil, fmt.Errorf("Failed GetPropertySSID(): %s", err.Error())
		}
		frequencyMHz, err := accessPoint.GetPropertyFrequency()
		if err != nil {
			return nil, fmt.Errorf("Failed GetPropertyFrequency(): %s", err.Error())
		}
		strength, err := accessPoint.GetPropertyStrength()
		if err != nil {
			return nil, fmt.Errorf("Failed GetPropertyStrength(): %s", err.Error())
		}
		scanResults = append(
			scanResults,
			wifi.ScanResult{
				Iface:             iface,
				BSSID:             bssid,
				SSID:              ssid,
				FrequencyMHz:      frequencyMHz,
				SignalStrengthdBm: getdBm(strength),
			},
		)
	}
	return scanResults, nil
}
