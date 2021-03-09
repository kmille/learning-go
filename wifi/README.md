## wifi client in go 

A small tool that helps you connecting to a wifi network. Written in Go.  Speaks to wpa_supplicant via d-bus interface. For debugging, use d-feet. 

## Usage

```bash
kmille@linbox:wifi sudo make tests
[sudo] password for kmille: 
go test -v
Let's do some tests
=== RUN   TestRegisterWifiInterface
Added interface wlp3s0
--- PASS: TestRegisterWifiInterface (0.07s)
=== RUN   TestDbusRegisterWifiInterface
Removed interface wlp3s0
Added interface wlp3s0
Could not create new interface: wpa_supplicant couldn't grab this interface.
--- PASS: TestDbusRegisterWifiInterface (0.13s)
=== RUN   TestRemoveWifiInterface
Removed interface wlp3s0
Added interface wlp3s0
Removed interface wlp3s0
--- PASS: TestRemoveWifiInterface (0.16s)
=== RUN   TestWifiScan
Added interface wlp3s0
Triggered a wifi scan
 0 MAC: dedacted SSID=error451                       signal=-62 freq=5745 wpa=           wpa2=wpa-psk wpa-psk-sha256 sae     open=false
 1 MAC: dedacted SSID=FRITZ!Box 3490                 signal=-69 freq=5500 wpa=           wpa2=wpa-psk                        open=false
 2 MAC: dedacted SSID=error451                       signal=-44 freq=2472 wpa=           wpa2=wpa-psk wpa-psk-sha256 sae     open=false
 3 MAC: dedacted SSID=FRITZ!Box 3490                 signal=-58 freq=2437 wpa=           wpa2=wpa-psk                        open=false
 4 MAC: dedacted SSID=FRITZ!Box 7560 UV              signal=-65 freq=2437 wpa=           wpa2=wpa-psk                        open=false
 5 MAC: dedacted SSID=Tantrafreunde-eV               signal=-65 freq=2462 wpa=           wpa2=wpa-psk                        open=false
 6 MAC: dedacted SSID=FRITZ!Box 7560 UV              signal=-79 freq=5260 wpa=           wpa2=wpa-psk                        open=false
 7 MAC: dedacted SSID=o2-WLAN42                      signal=-70 freq=2432 wpa=           wpa2=wpa-psk                        open=false
 8 MAC: dedacted SSID=SaMiLi                         signal=-74 freq=2437 wpa=wpa-psk    wpa2=wpa-psk                        open=false
 9 MAC: dedacted SSID=TP-Link_8AE2                   signal=-83 freq=2422 wpa=           wpa2=wpa-psk                        open=false
10 MAC: dedacted SSID=FRITZ!Box 3490                 signal=-79 freq=2412 wpa=           wpa2=wpa-psk                        open=false
11 MAC: dedacted SSID=FRITZ!Box 7490 HM              signal=-79 freq=2462 wpa=           wpa2=wpa-psk                        open=false
12 MAC: dedacted SSID=WLAN-483944                    signal=-80 freq=2462 wpa=           wpa2=wpa-psk                        open=false
13 MAC: dedacted SSID=Brokkoli                       signal=-80 freq=2412 wpa=           wpa2=wpa-psk                        open=false
14 MAC: dedacted SSID=PORTOHANYS                     signal=-83 freq=2462 wpa=           wpa2=wpa-psk                        open=false
15 MAC: dedacted SSID=IRONMAN                        signal=-83 freq=2452 wpa=           wpa2=wpa-psk                        open=false
16 MAC: dedacted SSID=Telekom_FON                    signal=-80 freq=2462 wpa=           wpa2=                               open=true
17 MAC: dedacted SSID=rockrobo-vacuum-v1_miap9756    signal=-83 freq=2412 wpa=           wpa2=                               open=true
    wifi_test.go:87: Scanned SSID: "error451"
Removed interface wlp3s0
--- PASS: TestWifiScan (2.24s)
PASS
Let's clean up
ok      wifi    2.603s
kmille@linbox:wifi 


kmille@linbox:wifi sudo make      
go build wifi.go


kmille@linbox:wifi ./wifi -h
Usage: wifi client that can scan for wifi networks and connect to it (uses d-bus interface of wpa_supplicant) 
-i string     name of your wifi network device, e.g. wlan0
scan          scan for networks
connect       scan for networks and connect to one
kmille@linbox:wifi 


kmille@linbox:wifi sudo ./wifi -i wlp3s0 scan
Added interface wlp3s0
Triggered a wifi scan
 0 MAC: dedacted SSID=error451                       signal=-60 freq=5745 wpa=           wpa2=wpa-psk wpa-psk-sha256 sae     open=false
 1 MAC: dedacted SSID=FRITZ!Box 3490                 signal=-69 freq=5500 wpa=           wpa2=wpa-psk                        open=false
 2 MAC: dedacted SSID=error451                       signal=-44 freq=2472 wpa=           wpa2=wpa-psk wpa-psk-sha256 sae     open=false
 3 MAC: dedacted SSID=FRITZ!Box 3490                 signal=-59 freq=2437 wpa=           wpa2=wpa-psk                        open=false
 4 MAC: dedacted SSID=Tantrafreunde-eV               signal=-67 freq=2462 wpa=           wpa2=wpa-psk                        open=false
 5 MAC: dedacted SSID=FRITZ!Box 7560 UV              signal=-68 freq=2437 wpa=           wpa2=wpa-psk                        open=false
 6 MAC: dedacted SSID=o2-WLAN42                      signal=-79 freq=5180 wpa=           wpa2=wpa-psk                        open=false
 7 MAC: dedacted SSID=o2-WLAN42                      signal=-70 freq=2432 wpa=           wpa2=wpa-psk                        open=false
 8 MAC: dedacted SSID=SaMiLi                         signal=-71 freq=2437 wpa=wpa-psk    wpa2=wpa-psk                        open=false
 9 MAC: dedacted SSID=Brokkoli                       signal=-79 freq=2412 wpa=           wpa2=wpa-psk                        open=false
10 MAC: dedacted SSID=FRITZ!Box 3490                 signal=-79 freq=2412 wpa=           wpa2=wpa-psk                        open=false
11 MAC: dedacted SSID=FRITZ!Box 4040 OW              signal=-79 freq=2437 wpa=           wpa2=wpa-psk                        open=false
12 MAC: dedacted SSID=KOLDNETT                       signal=-79 freq=2437 wpa=           wpa2=wpa-psk                        open=false
13 MAC: dedacted SSID=WLAN-483944                    signal=-82 freq=2462 wpa=           wpa2=wpa-psk                        open=false
14 MAC: dedacted SSID=Skynet                         signal=-83 freq=2437 wpa=           wpa2=wpa-psk                        open=false
15 MAC: dedacted SSID=o2-WLAN22                      signal=-83 freq=2437 wpa=           wpa2=wpa-psk                        open=false
16 MAC: dedacted SSID=IRONMAN                        signal=-83 freq=2452 wpa=           wpa2=wpa-psk                        open=false
17 MAC: dedacted SSID=ᛇ                              signal=-83 freq=2462 wpa=           wpa2=wpa-psk                        open=false
18 MAC: dedacted SSID=Telekom_FON                    signal=-80 freq=2462 wpa=           wpa2=                               open=true
19 MAC: dedacted SSID=Telekom_FON                    signal=-83 freq=2412 wpa=           wpa2=                               open=true
^CError: interrupt
Cleaning up
Removed interface wlp3s0
kmille@linbox:wifi 



kmille@linbox:wifi sudo ./wifi -i wlp3s0 connect
Added interface wlp3s0
Triggered a wifi scan
 0 MAC: dedacted SSID=error451                       signal=-62 freq=5745 wpa=           wpa2=wpa-psk wpa-psk-sha256 sae     open=false
 1 MAC: dedacted SSID=FRITZ!Box 3490                 signal=-71 freq=5500 wpa=           wpa2=wpa-psk                        open=false
 2 MAC: dedacted SSID=FRITZ!Box 3490                 signal=-58 freq=2437 wpa=           wpa2=wpa-psk                        open=false
 3 MAC: dedacted SSID=Tantrafreunde-eV               signal=-63 freq=2462 wpa=           wpa2=wpa-psk                        open=false
 4 MAC: dedacted SSID=FRITZ!Box 7560 UV              signal=-65 freq=2437 wpa=           wpa2=wpa-psk                        open=false
 5 MAC: dedacted SSID=SaMiLi                         signal=-68 freq=2437 wpa=wpa-psk    wpa2=wpa-psk                        open=false
 6 MAC: dedacted SSID=FRITZ!Box 7560 UV              signal=-79 freq=5260 wpa=           wpa2=wpa-psk                        open=false
 7 MAC: dedacted SSID=o2-WLAN42                      signal=-71 freq=2432 wpa=           wpa2=wpa-psk                        open=false
 8 MAC: dedacted SSID=TP-Link_8AE2                   signal=-82 freq=2422 wpa=           wpa2=wpa-psk                        open=false
 9 MAC: dedacted SSID=WLAN-483944                    signal=-78 freq=2462 wpa=           wpa2=wpa-psk                        open=false
10 MAC: dedacted SSID=FRITZ!Box 3490                 signal=-81 freq=2412 wpa=           wpa2=wpa-psk                        open=false
11 MAC: dedacted SSID=Brokkoli                       signal=-83 freq=2412 wpa=           wpa2=wpa-psk                        open=false
12 MAC: dedacted SSID=IRONMAN                        signal=-83 freq=2452 wpa=           wpa2=wpa-psk                        open=false
13 MAC: dedacted SSID=ᛇ                              signal=-83 freq=2462 wpa=           wpa2=wpa-psk                        open=false
14 MAC: dedacted SSID=Telekom_FON                    signal=-76 freq=2462 wpa=           wpa2=                               open=true
Which wifi do you want to connect? 0
Please enter the password for wifi "error451"
redacted
Added new wifi configuration
Successfully connected to "error451"
NOTE: the connection to the new wifi will be dropped if this program terminates.
Do you want to save the network profile for error451 in /etc/wpa_supplicant/wpa_supplicant-wlp3s0.conf?
y
Wrote config
Waiting for termination ...
```

## docs
- https://w1.fi/wpa_supplicant/devel/dbus.html
- https://pkg.go.dev/github.com/godbus/dbus/v5
