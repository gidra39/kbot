APP=$(shell basename $(shell git remote get-url origin))
REGISTRY=ghcr.io
USER=gidra39
VERSION= $(shell git describe --tags --abbrev=0)-$(shell git rev-parse --short HEAD)
TARGETOS=linux 
TARGETARCH=amd64
format:
	gofmt -s -w ./

lint:
	golint

test:
	go test -v

get:
	go get

linux:
	${MAKE} build TARGETOS=linux TARGETARCH=amd64

macOS:
	${MAKE} build TARGETOS=darwin TARGETARCH=amd64

windows:
	${MAKE} build TARGETOS=windows TARGETARCH=amd64 CGO_ENABLED=1

arm:
	${MAKE} build TARGETOS=linux TARGETARCH=arm64

build: format get
	CGO_ENABLED=0 GOOS=${TARGETOS} GOARCH=${TARGETARCH} go build -v -o kbot -ldflags "-X="github.com/gidra39/kbot/cmd.appVersion=${VERSION}

image:
	docker build . -t ${REGISTRY}/${USER}/${APP}:${VERSION} --build-arg CGO_ENABLED=${CGO_ENABLED} --build-arg TARGETARCH=${TARGETARCH} --build-arg TARGETOS=${TARGETOS}

push:
	docker push ${REGISTRY}/${USER}/${APP}:${VERSION}

clean: 
	docker rmi ${REGISTRY}/${USER}/${APP}:${VERSION}
	rm -rf kbot
