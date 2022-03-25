SHELL = /usr/bin/env bash

# File containing main function
SRC = ./...

# Build directory
BUILD_DIR = bin

# Destination directory
DEST_DIR = /usr/local/bin

.PHONY: build
build: clean
	go build -ldflags=-w -o ${BUILD_DIR}/ ${SRC}

.PHONY: clean
clean:
	rm -rf ${BUILD_DIR}/

.PHONY: install
install:
	systemctl stop rustymon-server.service
	cp -r ${BUILD_DIR}/* ${DEST_DIR}/
	mkdir -p /etc/rustymon-server/
	cp example.config.toml /etc/rustymon-server/
	cp rustymon-server.service /usr/lib/systemd/system/
	if [ -L /etc/systemd/system/multi-user.target.wants/rustymon-server.service ] ; then \
		if [ -e /etc/systemd/system/multi-user.target.wants/rustymon-server.service ]; then \
			echo "Service file is already linked properly"; \
		else \
			rm /etc/systemd/system/multi-user.target.wants/rustymon-server.service; \
			ln -s /usr/lib/systemd/system/rustymon-server.service /etc/systemd/system/multi-user.target.wants/; \
		fi \
	else \
		ln -s /usr/lib/systemd/system/rustymon-server.service /etc/systemd/system/multi-user.target.wants/; \
	fi
	systemctl daemon-reload

.PHONY: uninstall
uninstall:
	rm ${DEST_DIR}/rustymon-server ${DEST_DIR}/rustymon
	rm -rf /etc/rustymon-server/
