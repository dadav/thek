.PHONY: build
build:
	go build -o thek .

.PHONY: buildarm64
buildarm64:
	GOOS=linux GOARCH=arm64 go build -o thek_arm64

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
