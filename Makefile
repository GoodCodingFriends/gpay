SHELL := /bin/bash

dep:
ifeq ($(shell which dep 2>/dev/null),)
	curl https://raw.githubusercontent.com/golang/dep/master/install.sh | sh
endif

deps: dep
	@dep ensure

build: deps
	@go build 

test: gotest golint govet

gotest:
	@go test -race -v $(shell go list ./...)

golint:
	@# TODO: refactor
	@$(eval out := $(shell golint $(shell go list ./...) | grep -v 'have comment'))
	@test -z "$(out)" >/dev/null 2>&1 || (echo -e $(out) && false)

govet:
	@go vet $(shell go list ./... | grep -v repositorytest)

coverage: 
	@go tool cover -html=coverage.out

migrate:
	@mysql -h $(REPOSITORY_MYSQL_ADDRESS) -u $(REPOSITORY_MYSQL_USER) < database/schema.sql

.PHONY: dep deps build test gotest golint govet coverage
