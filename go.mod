module github.com/fornellas/wifi_exporter

go 1.17

replace pifke.org/wpasupplicant => ./golang-wpasupplicant

require (
	github.com/Wifx/gonetworkmanager v0.4.0
	gopkg.in/alecthomas/kingpin.v2 v2.2.6
	pifke.org/wpasupplicant v0.0.0-00010101000000-000000000000
)

require (
	github.com/alecthomas/template v0.0.0-20190718012654-fb15b899a751 // indirect
	github.com/alecthomas/units v0.0.0-20211218093645-b94a6e3cc137 // indirect
	github.com/godbus/dbus/v5 v5.0.2 // indirect
	github.com/stretchr/testify v1.8.0 // indirect
)
