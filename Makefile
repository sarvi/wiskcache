

SRCS:=*.go **/*.go

.PHONY: clean
clean:
	go clean
	rm wiskcache

./wiskcache: *.go
	go build

.PHONY: install
install: ./wiskcache
	install -D ./wiskcache $(ROOT)/bin/
	install -D ./wisk/config/wiskcache_config.yaml $(ROOT)/config

.PHONY: all
all: wiskcache
