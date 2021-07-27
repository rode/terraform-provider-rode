MAKEFLAGS += --silent

.PHONY: build install example

VERSION=0.0.1
GOOS=$(shell go env GOOS)
GOARCH=$(shell go env GOARCH)
REGISTRY=registry.terraform.io
NAMESPACE=rode
NAME=rode

build:
	go build -o terraform-provider-rode

install: build
	mkdir -p ~/.terraform.d/plugins/${REGISTRY}/${NAMESPACE}/${NAME}/${VERSION}/${GOOS}_${GOARCH}/
	mv terraform-provider-rode ~/.terraform.d/plugins/${REGISTRY}/${NAMESPACE}/${NAME}/${VERSION}/${GOOS}_${GOARCH}/

example: build
	terraform -chdir=example apply
