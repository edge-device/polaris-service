SHELL:=/bin/bash
PROJECT_NAME=polaris
GO_BUILD_ENV=CGO_ENABLED=0 GOOS=linux GOARCH=amd64
GO_FILES=$(shell go list ./... | grep -v /vendor/)

BUILD_VERSION=$(shell cat VERSION)
BUILD_TAG=$(BUILD_VERSION)
DOCKER_IMAGE=$(PROJECT_NAME):$(BUILD_TAG)

.SILENT:

all: fmt vet install test

build:
	$(GO_BUILD_ENV) go build -v -o $(PROJECT_NAME) *.go

install:
	$(GO_BUILD_ENV) go install

vet:
	go vet $(GO_FILES)

fmt:
	go fmt $(GO_FILES)

test:
	go test $(GO_FILES) -cover

docker: build
	docker build -t $(DOCKER_IMAGE) .;\
        rm -f $(PROJECT_NAME).bin 2> /dev/null; \