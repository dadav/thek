.PHONY: build
build:
	go build -o thek .

.PHONY: install
install:
	cp thek /usr/local/bin/
	mkdir -p /etc/thek
	cp config.yaml /etc/thek/
	cp thek.service /etc/systemd/system/
	systemctl daemon-reload

.PHONY: clean
clean:
	rm thek
