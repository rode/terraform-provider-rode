MAKEFLAGS += --silent

.PHONY: build install example fmtcheck tfcheck fmt testacc

VERSION=0.0.1
GOOS=$(shell go env GOOS)
GOARCH=$(shell go env GOARCH)
REGISTRY=registry.terraform.io
NAMESPACE=rode
NAME=rode

testacc:
	TF_ACC=1 RODE_HOST=localhost:50051 go test ./...

build:
	go build -o terraform-provider-rode

install: build
	mkdir -p ~/.terraform.d/plugins/${REGISTRY}/${NAMESPACE}/${NAME}/${VERSION}/${GOOS}_${GOARCH}/
	mv terraform-provider-rode ~/.terraform.d/plugins/${REGISTRY}/${NAMESPACE}/${NAME}/${VERSION}/${GOOS}_${GOARCH}/

example: build
	terraform -chdir=example apply

fmtcheck:
	exit $(shell gofmt -s -l . | wc -l)

tfcheck:
	terraform fmt -recursive -check

fmt:
	gofmt -s -w .

test: fmtcheck tfcheck
	go vet ./...
	go test -v ./...
