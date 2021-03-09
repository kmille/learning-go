package main

import (
	"fmt"
	"testing"

	"github.com/godbus/dbus/v5"
)

var conn *dbus.Conn

const ifaceTest string = "wlp3s0"

func TestMain(m *testing.M) {
	fmt.Println("Let's do some tests")
	conn = connect()
	checkRootPrivileges()
	m.Run()
	fmt.Println("Let's clean up")
	conn.Close()
}

func TestRegisterWifiInterface(t *testing.T) {
	_, err := getWifiInterfaceName(conn, "this does not exist")
	if err == nil {
		t.Fatal("This interface should not exist", err)
	}

	err = registerWifiInterface(conn, ifaceTest)

	if err != nil {
		t.Fatal("Could not register interface for testing", err)
	}
	_, err = getWifiInterfaceName(conn, ifaceTest)
	if err != nil {
		t.Fatal("Could not get interface name", err)
	}
}

func TestDbusRegisterWifiInterface(t *testing.T) {
	err := removeWifiInterface(conn, ifaceTest)
	if err != nil {
		t.Fatal("Could not remove interface for testing", err)
	}
	err = registerWifiInterface(conn, ifaceTest)
	if err != nil {
		t.Fatal("Could not register interface for testing", err)
	}

	err = registerWifiInterface(conn, "this does not exist")
	if err == nil {
		t.Fatal("Registering an unknown interface should fail", err)
	}
}

func TestRemoveWifiInterface(t *testing.T) {
	err := removeWifiInterface(conn, ifaceTest)
	if err != nil {
		t.Fatal("Could not remove interface", err)
	}
	if err = registerWifiInterface(conn, ifaceTest); err != nil {
		t.Fatal("Could not register interface", err)
	}
	err = removeWifiInterface(conn, ifaceTest)
	if err != nil {
		t.Fatal("Could not remove interface", err)
	}
}

func TestWifiScan(t *testing.T) {
	err := registerWifiInterface(conn, ifaceTest)
	if err != nil {
		t.Fatal("Could not register interface for testing", err)

	}
	scanWifiNetworks(conn, ifaceTest, true)
	scannedWifis := handleScanComplete(conn, ifaceTest, true)
	if err != nil {
		t.Fatal("Could not get scanned wifi networks", err)
	}
	variant, err := getWifiProperty(conn, dbus.ObjectPath(scannedWifis[0].dbusIdentifier), "SSID")
	if err != nil {
		t.Fatal("Could not read SSID")
	}

	ssid := variant.Value().([]uint8)
	t.Logf("Scanned SSID: %q", ssid)

	err = removeWifiInterface(conn, ifaceTest)
	if err != nil {
		t.Fatal("Could not remove interface", err)
	}
	// func getWifiProperty(conn *dbus.Conn, objectPath dbus.ObjectPath, property string) (dbus.Variant, error) {
}
