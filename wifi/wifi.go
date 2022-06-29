package wifi

import (
	"net"
	"time"
)

type Iface string

func NewIface(name string) Iface {
	return Iface(name)
}

func (i Iface) Name() string {
	return string(i)
}

type ScanResult struct {
	// Iface is the interface name
	Iface Iface
	// BSSID is the MAC address of the BSS.
	BSSID net.HardwareAddr
	// SSID is the SSID of the BSS.
	SSID string
	// FrequencyMHz is the frequency, in MHz, of the BSS.
	FrequencyMHz uint32
	// Received signal strength, in dB, of the BSS.
	SignalStrengthdBm int
}

func (s *ScanResult) FrequencyBand() string {
	if s.FrequencyMHz >= 2412 && s.FrequencyMHz <= 2484 {
		return "2.4GHz"
	}
	if s.FrequencyMHz >= 5035 && s.FrequencyMHz <= 5980 {
		return "5GHz"
	}
	return "Unknown"
}

func (s *ScanResult) Channel() int {
	// 2.4GHz https://en.wikipedia.org/wiki/List_of_WLAN_channels#2.4_GHz_(802.11b/g/n/ax)
	if s.FrequencyMHz == 2412 {
		return 1
	}
	if s.FrequencyMHz == 2417 {
		return 2
	}
	if s.FrequencyMHz == 2422 {
		return 3
	}
	if s.FrequencyMHz == 2427 {
		return 4
	}
	if s.FrequencyMHz == 2432 {
		return 5
	}
	if s.FrequencyMHz == 2437 {
		return 6
	}
	if s.FrequencyMHz == 2442 {
		return 7
	}
	if s.FrequencyMHz == 2447 {
		return 8
	}
	if s.FrequencyMHz == 2452 {
		return 9
	}
	if s.FrequencyMHz == 2457 {
		return 10
	}
	if s.FrequencyMHz == 2462 {
		return 11
	}
	if s.FrequencyMHz == 2467 {
		return 12
	}
	if s.FrequencyMHz == 2472 {
		return 13
	}
	if s.FrequencyMHz == 2484 {
		return 14
	}

	// 5GHz https://en.wikipedia.org/wiki/List_of_WLAN_channels#5_GHz_(802.11a/h/j/n/ac/ax)
	if s.FrequencyMHz == 5035 {
		return 7
	}
	if s.FrequencyMHz == 5040 {
		return 8
	}
	if s.FrequencyMHz == 5045 {
		return 9
	}
	if s.FrequencyMHz == 5055 {
		return 11
	}
	if s.FrequencyMHz == 5060 {
		return 12
	}
	if s.FrequencyMHz == 5080 {
		return 16
	}
	if s.FrequencyMHz == 5160 {
		return 32
	}
	if s.FrequencyMHz == 5170 {
		return 34
	}
	if s.FrequencyMHz == 5180 {
		return 36
	}
	if s.FrequencyMHz == 5190 {
		return 38
	}
	if s.FrequencyMHz == 5200 {
		return 40
	}
	if s.FrequencyMHz == 5210 {
		return 42
	}
	if s.FrequencyMHz == 5220 {
		return 44
	}
	if s.FrequencyMHz == 5230 {
		return 46
	}
	if s.FrequencyMHz == 5240 {
		return 48
	}
	if s.FrequencyMHz == 5250 {
		return 50
	}
	if s.FrequencyMHz == 5260 {
		return 52
	}
	if s.FrequencyMHz == 5270 {
		return 54
	}
	if s.FrequencyMHz == 5280 {
		return 56
	}
	if s.FrequencyMHz == 5290 {
		return 58
	}
	if s.FrequencyMHz == 5300 {
		return 60
	}
	if s.FrequencyMHz == 5310 {
		return 62
	}
	if s.FrequencyMHz == 5320 {
		return 64
	}
	if s.FrequencyMHz == 5340 {
		return 68
	}
	if s.FrequencyMHz == 5480 {
		return 96
	}
	if s.FrequencyMHz == 5500 {
		return 100
	}
	if s.FrequencyMHz == 5510 {
		return 102
	}
	if s.FrequencyMHz == 5520 {
		return 104
	}
	if s.FrequencyMHz == 5530 {
		return 106
	}
	if s.FrequencyMHz == 5540 {
		return 108
	}
	if s.FrequencyMHz == 5550 {
		return 110
	}
	if s.FrequencyMHz == 5560 {
		return 112
	}
	if s.FrequencyMHz == 5570 {
		return 114
	}
	if s.FrequencyMHz == 5580 {
		return 116
	}
	if s.FrequencyMHz == 5590 {
		return 118
	}
	if s.FrequencyMHz == 5600 {
		return 120
	}
	if s.FrequencyMHz == 5610 {
		return 122
	}
	if s.FrequencyMHz == 5620 {
		return 124
	}
	if s.FrequencyMHz == 5630 {
		return 126
	}
	if s.FrequencyMHz == 5640 {
		return 128
	}
	if s.FrequencyMHz == 5660 {
		return 132
	}
	if s.FrequencyMHz == 5670 {
		return 134
	}
	if s.FrequencyMHz == 5680 {
		return 136
	}
	if s.FrequencyMHz == 5690 {
		return 138
	}
	if s.FrequencyMHz == 5700 {
		return 140
	}
	if s.FrequencyMHz == 5710 {
		return 142
	}
	if s.FrequencyMHz == 5720 {
		return 144
	}
	if s.FrequencyMHz == 5745 {
		return 149
	}
	if s.FrequencyMHz == 5755 {
		return 151
	}
	if s.FrequencyMHz == 5765 {
		return 153
	}
	if s.FrequencyMHz == 5775 {
		return 155
	}
	if s.FrequencyMHz == 5785 {
		return 157
	}
	if s.FrequencyMHz == 5795 {
		return 159
	}
	if s.FrequencyMHz == 5805 {
		return 161
	}
	if s.FrequencyMHz == 5815 {
		return 163
	}
	if s.FrequencyMHz == 5825 {
		return 165
	}
	if s.FrequencyMHz == 5835 {
		return 167
	}
	if s.FrequencyMHz == 5845 {
		return 169
	}
	if s.FrequencyMHz == 5855 {
		return 171
	}
	if s.FrequencyMHz == 5865 {
		return 173
	}
	if s.FrequencyMHz == 5875 {
		return 175
	}
	if s.FrequencyMHz == 5885 {
		return 177
	}
	if s.FrequencyMHz == 5910 {
		return 182
	}
	if s.FrequencyMHz == 5915 {
		return 183
	}
	if s.FrequencyMHz == 5920 {
		return 184
	}
	if s.FrequencyMHz == 5935 {
		return 187
	}
	if s.FrequencyMHz == 5940 {
		return 188
	}
	if s.FrequencyMHz == 5945 {
		return 189
	}
	if s.FrequencyMHz == 5960 {
		return 192
	}
	if s.FrequencyMHz == 5980 {
		return 196
	}

	return 0
}

type Client interface {
	GetInterfaces() (ifaces []Iface, err error)
	Scan(iface Iface, timeout time.Duration) ([]ScanResult, error)
}
