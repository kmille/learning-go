package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/godbus/dbus/v5"
)

const busName string = "fi.w1.wpa_supplicant1"
const dbusInterfacePrefix string = "fi.w1.wpa_supplicant1"
const dbusObjectPath = dbus.ObjectPath("/fi/w1/wpa_supplicant1")
const wifiScanInterval time.Duration = 5

type wifi struct {
	dbusIdentifier string
	ssid           []uint8
	signal         int16
	open           bool
	wpa            string
	wpa2           string
	macAddress     []byte
	frequency      uint16
}

func (w *wifi) String() string {
	return fmt.Sprintf("MAC:%s SSID=%-30s signal=%d freq=%d wpa=%-10s wpa2=%-30s open=%t", hexEncodeMacAddress(w.macAddress),
		string(w.ssid), w.signal, w.frequency, w.wpa, w.wpa2, w.open)
}

func hexEncodeMacAddress(mac []byte) string {
	var macEncoded string
	for _, b := range mac {
		macEncoded += fmt.Sprintf("%02x:", b)
	}
	// ignore last :
	return macEncoded[:len(macEncoded)-1]
}

func newWifi(conn *dbus.Conn, dbusIdentifier dbus.ObjectPath) *wifi {
	b := wifi{dbusIdentifier: string(dbusIdentifier)}
	variant, err := getWifiProperty(conn, dbusIdentifier, "SSID")
	if err != nil {
		fmt.Println("TODO error handling 1")
		return nil
	}
	b.ssid = variant.Value().([]uint8)
	variant, err = getWifiProperty(conn, dbusIdentifier, "Signal")
	if err != nil {
		fmt.Println("TODO error handling 2")
		return nil
	}
	b.signal = variant.Value().(int16)

	variant, err = getWifiProperty(conn, dbusIdentifier, "Frequency")
	if err != nil {
		fmt.Println("TODO error handling 2")
		return nil
	}
	b.frequency = variant.Value().(uint16)

	variant, err = getWifiProperty(conn, dbusIdentifier, "WPA")
	if err != nil {
		fmt.Println("TODO error handling 2")
		return nil
	}
	wpa := variant.Value().(map[string]dbus.Variant)
	wpaKeyMgmt := wpa["KeyMgmt"].Value().([]string)
	b.wpa = strings.Join(wpaKeyMgmt, " ")

	variant, err = getWifiProperty(conn, dbusIdentifier, "RSN")
	if err != nil {
		fmt.Println("TODO error handling 3")
		return nil
	}
	wpa2 := variant.Value().(map[string]dbus.Variant)
	wpa2KeyMgmt := wpa2["KeyMgmt"].Value().([]string)
	b.wpa2 = strings.Join(wpa2KeyMgmt, " ")

	variant, err = getWifiProperty(conn, dbusIdentifier, "BSSID")
	if err != nil {
		fmt.Println("TODO error handling 4")
		return nil
	}
	b.macAddress = variant.Value().([]byte)

	variant, err = getWifiProperty(conn, dbusIdentifier, "Privacy")
	if err != nil {
		fmt.Println("TODO error handling 4")
		return nil
	}
	b.open = !variant.Value().(bool)

	return &b
}

func getWifiProperty(conn *dbus.Conn, objectPath dbus.ObjectPath, property string) (dbus.Variant, error) {
	variant, err := conn.Object(busName, objectPath).GetProperty(busName + ".BSS." + property)
	if err != nil {
		fmt.Printf("Failed to read wifi property %q: %s\n", property, err)
		return dbus.MakeVariant(""), err
	}
	return variant, nil
}

func registerWifiInterface(conn *dbus.Conn, iface string) error {
	argument := make(map[string]dbus.Variant)
	argument["Ifname"] = dbus.MakeVariant(iface)
	obj := conn.Object(busName, dbusObjectPath)
	call := obj.Call(dbusInterfacePrefix+".CreateInterface", 0, argument)
	if call.Err != nil {
		if strings.Contains(call.Err.Error(), "wpa_supplicant already controls this interface") {
			fmt.Println("Device already added:", call.Err)
			return nil
		} else {
			fmt.Println("Could not create new interface:", call.Err)
			return call.Err
		}
	}
	fmt.Println("Added interface", iface)
	return nil
}

func removeWifiInterface(conn *dbus.Conn, iface string) error {
	obj := conn.Object(busName, dbusObjectPath)
	name, err := dbusGetWifiInterfaceName(conn, iface)
	if err != nil {
		// interface was removed previously
		return nil
	}
	call := obj.Call(dbusInterfacePrefix+".RemoveInterface", 0, dbus.ObjectPath(name))
	if call.Err != nil {
		fmt.Println("Could not remove interface:", call.Err)
		return call.Err
	}
	fmt.Println("Removed interface", iface)
	return nil
}

func dbusGetWifiInterfaceName(conn *dbus.Conn, iface string) (interfacePath string, err error) {
	obj := conn.Object(busName, dbusObjectPath)
	err = obj.Call(dbusInterfacePrefix+".GetInterface", 0, iface).Store(&interfacePath)
	return
}

func scanWifiNetworks(conn *dbus.Conn, iface string, runOnce bool) {
	for {
		interfacePath, err := dbusGetWifiInterfaceName(conn, iface)
		if err != nil {
			cleanup(conn, iface, err)
		}
		obj := conn.Object("fi.w1.wpa_supplicant1", dbus.ObjectPath(interfacePath))
		argument := make(map[string]dbus.Variant)
		argument["Type"] = dbus.MakeVariant("active")
		call := obj.Call("fi.w1.wpa_supplicant1.Interface.Scan", 0, argument)
		if call.Err != nil {
			fmt.Println("Could not scan wifi signals:", call.Err)
			cleanup(conn, iface, err)
		}
		fmt.Println("Triggered a wifi scan")

		if runOnce {
			return
		} else {
			time.Sleep(wifiScanInterval * time.Second)
		}
	}
}

func getScannedWifiNetworks(conn *dbus.Conn, iface string) ([]*wifi, error) {
	var wifis []*wifi
	interfaceName, err := dbusGetWifiInterfaceName(conn, iface)
	if err != nil {
		return []*wifi{}, err
	}
	variant, err := conn.Object(busName, dbus.ObjectPath(interfaceName)).GetProperty(busName + ".Interface.BSSs")
	if err != nil {
		fmt.Println("Failed to get scanned wifi networks:", err)
		return []*wifi{}, err
	}
	wifiNetworkIdentifiers := variant.Value().([]dbus.ObjectPath)
	for _, wifiIdentifier := range wifiNetworkIdentifiers {
		wifis = append(wifis, newWifi(conn, wifiIdentifier))
	}
	return wifis, nil
}

func printFoundWifis(conn *dbus.Conn, iface string, runOnce bool) ([]*wifi, error) {
	wifis, err := getScannedWifiNetworks(conn, iface)
	if err != nil {
		return []*wifi{}, err
	}
	if !runOnce {
		// clear screen
		fmt.Print("\033[H\033[2J")
	}
	for i, wifi := range wifis {
		fmt.Printf("%2d %s\n", i, wifi)
	}
	return wifis, nil
}

func handleScanComplete(conn *dbus.Conn, iface string, runOnce bool) []*wifi {
	// https://github.com/godbus/dbus/blob/master/_examples/signal.go
	// https://dbus.freedesktop.org/doc/dbus-specification.html#message-bus-routing-match-rules
	err := conn.AddMatchSignal(dbus.WithMatchMember("ScanDone"))
	if err != nil {
		cleanup(conn, iface, err)
	}
	c := make(chan *dbus.Signal, 1)
	conn.Signal(c)
	for range c {
		wifis, err := printFoundWifis(conn, iface, runOnce)
		if err != nil {
			cleanup(conn, iface, err)
		}
		if runOnce {
			return wifis
		}
	}
	return nil
}

func connectNewWifiNetwork(conn *dbus.Conn, iface string, wifi *wifi) {
	var password string
	argument := make(map[string]dbus.Variant)
	argument["ssid"] = dbus.MakeVariant(wifi.ssid)
	if wifi.open {
		password = "NONE"
	} else {
		fmt.Println("Please enter the wifi password")
		reader := bufio.NewReader(os.Stdin)
		text, _ := reader.ReadString('\n')
		// TODO: calculate hash
		// argument["psk"] = dbus.MakeVariant("AiNuf7xeBeaJei0a")
		password = strings.TrimRight(text, "\n")
	}
	argument["psk"] = dbus.MakeVariant(password)
	if strings.Contains(wifi.wpa2, "sae") {
		// handle WPA3
		argument["key_mgmt"] = dbus.MakeVariant("SAE")
		argument["ieee80211w"] = dbus.MakeVariant(2)
	}

	name, err := dbusGetWifiInterfaceName(conn, iface)
	if err != nil {
		cleanup(conn, iface, err)
	}

	obj := conn.Object(busName, dbus.ObjectPath(name))
	call := obj.Call(dbusInterfacePrefix+".Interface.AddNetwork", 0, argument)
	if call.Err != nil {
		cleanup(conn, iface, call.Err)
	}
	fmt.Println("Added new wifi configuration")

	call = obj.Call(dbusInterfacePrefix+".Interface.SelectNetwork", 0, call.Body[0])
	if call.Err != nil {
		cleanup(conn, iface, call.Err)
	}
	fmt.Printf("Successfully switched to the new network %q\n", wifi.ssid)

	fmt.Println("NOTE: the connection to the new wifi will be dropped if this program terminates.")
	wpaSupplicantConfig := fmt.Sprintf("/etc/wpa_supplicant/wpa_supplicant-%s.conf", iface)
	fmt.Printf("Do you want to save the network profile for %s in %s?\n", wifi.ssid, wpaSupplicantConfig)

	// TODO: check if this file exists

	reader := bufio.NewReader(os.Stdin)
	text, _ := reader.ReadString('\n')
	if strings.Contains(text, "y") {
		f, err := os.OpenFile(wpaSupplicantConfig, os.O_APPEND|os.O_WRONLY, 0600)
		if err != nil {
			fmt.Println("Error writing config file:", err)
		} else {
			wpaSupplicantProfile := fmt.Sprintf("\nnetwork={         \nssid=%s         \npsk=\"%s\"\n", wifi.ssid, password)
			fmt.Println(wpaSupplicantProfile)
			_, err = f.WriteString(wpaSupplicantProfile)
			if err != nil {
				fmt.Println("Error writing config file:", err)
			} else {
				fmt.Println("Wrote config")
			}
		}
		f.Close()
	}

	// TODO: das ist mist hier. danach verschwindets wieder
	fmt.Println("Waiting ...")
	time.Sleep(50 * time.Second)
}

func checkRootPrivileges() {
	if os.Getuid() != 0 {
		fmt.Println("You have to be root")
		os.Exit(1)
	}
}

func connect() *dbus.Conn {
	conn, err := dbus.SystemBus()
	if err != nil {
		fmt.Println("Failed to connect to session bus:", err)
		os.Exit(1)
	}
	return conn
}

func cleanup(conn *dbus.Conn, iface string, msg interface{}) {
	if msg != nil {
		fmt.Println("Error:", msg)
	}
	fmt.Println("Cleaning up")
	err := removeWifiInterface(conn, iface)
	if err != nil {
		fmt.Println("Error during cleanup:", err)
	}
	err = conn.Close()
	if err != nil {
		fmt.Println("Error during cleanup:", err)
	}
	if msg != nil || err != nil {
		os.Exit(1)
	}
	os.Exit(0)
}

func main() {
	var iface string
	flag.Usage = func() {
		fmt.Println("Usage:\nTODO")
	}
	flag.StringVar(&iface, "i", "", "name of your wifi network device, e.g. wlan0")
	flag.Parse()
	if len(iface) == 0 || len(flag.Args()) != 1 {
		flag.Usage()
		os.Exit(1)
	}

	checkRootPrivileges()
	conn := connect()
	if err := registerWifiInterface(conn, iface); err != nil {
		cleanup(conn, iface, err)
	}

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		sig := <-sigs
		cleanup(conn, iface, sig)
	}()

	done := make(chan bool, 1)
	if flag.Args()[0] == "scan" {
		// both go routines are stopped via ctrl-c
		go handleScanComplete(conn, iface, false)
		go scanWifiNetworks(conn, iface, false)
	}

	if flag.Args()[0] == "connect" {
		scanWifiNetworks(conn, iface, true)
		wifis := handleScanComplete(conn, iface, true)
		reader := bufio.NewReader(os.Stdin)
		fmt.Println("Which wifi do you want to connect?")
		fmt.Print("Enter text: ")
		text, _ := reader.ReadString('\n')
		wifiIndex, err := strconv.Atoi(strings.TrimRight(text, "\n"))
		if err != nil {
			cleanup(conn, iface, err)
		}
		fmt.Println("you choose", wifiIndex, wifis[wifiIndex].dbusIdentifier)
		connectNewWifiNetwork(conn, iface, wifis[wifiIndex])
		done <- true
	}

	<-done
	cleanup(conn, iface, nil)

}
