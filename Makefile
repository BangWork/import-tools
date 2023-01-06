GOOS=$(shell uname | tr '[:upper:]' '[:lower:]')
GOARCH=$(shell uname -m)
CGO_ENABLED=0
ENTRY=main.go

.PHYNO: all
all: build-web copy-dist build-api

.PHYNO: build-web
build-web:
	cd web && npm i && npm run build

copy-dist: clean-dist
	cp -apr web/dist serve/router/

.PHYNO: build-api
build-api:
	GOOS=$(GOOS) GOARCH=$(GOARCH) CGO_ENABLED=$(CGO_ENABLED) go build \
		-o bin/$(GOOS)/import-tools \
		-trimpath \
		$(ENTRY)

build-linux:
	GOOS=linux GOARCH=amd64 CGO_ENABLED=$(CGO_ENABLED) go build \
		-o bin/$(GOOS)/import-tools \
		-trimpath \
		$(ENTRY)

clean-dist:
	rm -rf serve/router/dist
	rm -rf web/dist
