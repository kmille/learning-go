all:
	go build wifi.go

tests:
	go test -v

clean:
	rm -rf wifi

install:
	sudo install -m 0755 wifi /usr/local/bin/wifi

uninstall:
	sudo rm -rf /usr/local/bin/wifi
