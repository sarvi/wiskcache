

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

.PHONY: all
all: wiskcache
