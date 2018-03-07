SHELL := /bin/bash

glide:
ifeq ($(shell which glide 2>/dev/null),)
	curl https://glide.sh/get | sh
endif

deps: glide
	glide install

build: deps
	go build 

test: gotest golint govet

gotest:
	go test -race -v $(shell glide novendor)

golint:
	golint $(shell glide novendor)

govet:
	go vet $(shell glide novendor)

coverage: 
	go tool cover -html=coverage.out

.PHONY: glide deps build test gotest golint govet coverage
