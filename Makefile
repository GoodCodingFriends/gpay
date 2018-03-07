SHELL := /bin/bash

glide:
ifeq ($(shell which glide 2>/dev/null),)
	@curl https://glide.sh/get | sh
endif

deps: glide
	@glide install

build: deps
	@go build 

test: gotest golint govet

gotest:
	@go test -race -v $(shell glide novendor)

golint:
	# TODO: refactor
	@$(eval out := $(shell golint $(shell glide novendor) | grep -v 'have comment'))
	@test -z "$(out)" >/dev/null 2>&1 || echo -e $(out) && false

govet:
	@go vet $(shell go list $(shell glide novendor) | grep -v repositorytest)

coverage: 
	@go tool cover -html=coverage.out

.PHONY: glide deps build test gotest golint govet coverage
