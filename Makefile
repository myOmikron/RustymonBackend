# File containing main function
CMD = cmd/rustymon-server/rustymon-server.go

# Output directory
OUT_DIR = bin/

.PHONY: build
build: clean
	go build -ldflags=-w -o ${OUT_DIR} ${CMD}

.PHONY: clean
clean:
	rm -rf ${OUT_DIR}

.PHONY: install
install:
	echo "No installation yet"

.PHONY: uninstall
uninstall:
	echo "No uninstallation yet"
