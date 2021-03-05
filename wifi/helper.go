package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/godbus/dbus/v5"
	"github.com/godbus/dbus/v5/introspect"
)

func eavesdrop() {
	conn, err := dbus.SessionBus()
	if err != nil {
		fmt.Println("Failed to connect to session bus:", err)
		os.Exit(1)
	}
	defer conn.Close()

	for _, v := range []string{"method_call", "method_return", "error", "signal"} {
		call := conn.BusObject().Call("org.freedesktop.DBus.AddMatch", 0,
			"eavesdrop='true',type='"+v+"'")
		if call.Err != nil {
			fmt.Println("Failed to add match", call.Err)
			os.Exit(1)
		}
		c := make(chan *dbus.Message, 10)
		conn.Eavesdrop(c)
		fmt.Println("Listening for everything")
		for v := range c {
			fmt.Println(v)
		}

	}
}

func listNames() {
	conn, err := dbus.SessionBus()
	if err != nil {
		fmt.Println("Failed to connect to session bus:", err)
		os.Exit(1)
	}
	defer conn.Close()

	var s []string
	err = conn.BusObject().Call("org.freedesktop.DBus.ListNames", 0).Store(&s)
	if err != nil {
		fmt.Println("Failed to get list of wned names:", err)
		os.Exit(1)
	}

	fmt.Println("Currently onwed names on the session bus:")
	for _, v := range s {
		fmt.Println(v)
	}
}

func introspectEndpoint() {
	conn, err := dbus.SessionBus()
	if err != nil {
		fmt.Println("Failed to connect to session bus:", err)
		os.Exit(1)
	}
	defer conn.Close()
	node, err := introspect.Call(conn.Object("org.freedesktop.DBus", "/org/freedesktopp/Dbus"))
	if err != nil {
		panic(err)
	}
	data, _ := json.MarshalIndent(node, "", "     ")
	fmt.Printf("%s\n", data)
}

func sendNotification() {
	conn, err := dbus.SessionBus()
	if err != nil {
		fmt.Println("Failed to connect to session bus:", err)
		os.Exit(1)
	}
	// api docs: https://developer.gnome.org/notification-spec/ and search for org.freedesktop.Notifications.Notify
	obj := conn.Object("org.freedesktop.Notifications", "/org/freedesktop/Notifications")
	call := obj.Call("org.freedesktop.Notifications.Notify", 0, "", uint32(0),
		"", "Test", "This is a test of the DBus binding for go.", []string{},
		map[string]dbus.Variant{}, int32(5000))
	if call.Err != nil {
		panic(call.Err)
	}
}
