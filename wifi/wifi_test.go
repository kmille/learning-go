package main

import (
	"fmt"
	"testing"

	"github.com/godbus/dbus/v5"
)

const iface string = "wlp3s0"

var conn *dbus.Conn

func init() {
	// defer conn.Close()
}
func TestMain(m *testing.M) {
	fmt.Println("Let's do some tests")
	conn = connect()
	checkRootPrivileges()
	m.Run()
	fmt.Println("Let's clean up")
	conn.Close()
}

func TestDbusGetWifiInterfaceName(t *testing.T) {
	_, err := dbusGetWifiInterfaceName(conn, "this does not exist")
	if err == nil {
		t.Fatal("This interface should not exist", err)
	}

	err = dbusRegisterWifiInterface(conn, iface)
	if err != nil {
		t.Fatal("Could not register interface for testing", err)
	}
	_, err = dbusGetWifiInterfaceName(conn, iface)
	if err != nil {
		t.Fatal("Could not get interface name", err)
	}
}

func TestDbusRegisterWifiInterface(t *testing.T) {
	err := dbusRemoveWifiInterface(conn, iface)
	if err != nil {
		t.Fatal("Could not remove interface for testing", err)
	}
	err = dbusRegisterWifiInterface(conn, iface)
	if err != nil {
		t.Fatal("Could not register interface for testing", err)
	}

	err = dbusRegisterWifiInterface(conn, "this does not exist")
	if err == nil {
		t.Fatal("Registering an unknown interface should fail", err)
	}
}

func TestDbusRemoveWifiInterface(t *testing.T) {
	err := dbusRemoveWifiInterface(conn, iface)
	if err != nil {
		t.Fatal("Could not remove interface", err)
	}
	if err = dbusRegisterWifiInterface(conn, iface); err != nil {
		t.Fatal("Could not register interface", err)
	}
	err = dbusRemoveWifiInterface(conn, iface)
	if err != nil {
		t.Fatal("Could not remove interface", err)
	}
}

/*
func TestDbusListRegisteredWifiInterfaces(t *testing.T) {
	_, err := dbusListRegisteredWifiInterfaces(conn)
	if err != nil {
		t.Fatal("Error listing interfaces", err)
	}
	// t.Log("Found interface(s):", len(interfaces))
}
*/
