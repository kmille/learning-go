all:
	go build timer.go

clean:
	rm -rf timer

install:
	sudo install -m 0755 timer /usr/local/bin/timer

uninstall:
	rm -rf /usr/local/bin/timer

list:
	@grep '^[^#[:space:]].*:' Makefile

