GOOS = $(shell uname | tr '[:upper:]' '[:lower:]')
GOARCH = $(shell uname -m)
CGO_ENABLED = 0
ENTRY = main.go

IMAGE_TAG = latest
IMAGE_NAME = ghcr.io/bangwork/import-tools

.PHYNO: all
all: build-web build-api

.PHYNO: build-web
build-web:
	cd web && npm i && npm run build

copy-dist:
	rm -rf serve/router/dist
	cp -r web/dist serve/router/

.PHYNO: build-api
build-api: copy-dist
	GOOS=$(GOOS) GOARCH=$(GOARCH) CGO_ENABLED=$(CGO_ENABLED) go build \
		-o bin/$(GOOS)/import-tools \
		-trimpath \
		$(ENTRY)

.PHYNO: build-linux
build-linux: copy-dist
	GOOS=linux GOARCH=amd64 CGO_ENABLED=$(CGO_ENABLED) go build \
		-o bin/linux/import-tools \
		-trimpath \
		$(ENTRY)

.PHYNO: clean-dist
clean-dist:
	rm -rf serve/router/dist
	rm -rf web/dist

.PHYNO: build-image
build-image:
	docker build -t $(IMAGE_NAME):$(IMAGE_TAG) .
	[ "$(IMAGE_TAG)" != "latest" ] && docker tag $(IMAGE_NAME):$(IMAGE_TAG) $(IMAGE_NAME):latest || true

.PHYNO: push-image
push-image:
	docker push $(IMAGE_NAME):$(IMAGE_TAG)

.PHYNO: package
package:
	tar zcf import-tools.tar.gz \
		LICENSE \
		README.md \
		scripts \
		bin
