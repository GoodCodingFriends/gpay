SHELL := /bin/bash

REGISTRY ?= "gpay-gacha"
GACHA_CLOUD_RUN_REGION ?= "us-central1"

build/gacha:
	@go build -o gacha ./cmd/gacha 

image/gacha:
	@echo "building image..."
	@echo "registry: $(REGISTRY)"
	@docker build -t $(REGISTRY) -f ./cmd/gacha/Dockerfile .

remote/image/gacha: image/gacha
	@docker push $(REGISTRY)

deploy/gacha: remote/image/gacha
	@gcloud beta run deploy \
		--allow-unauthenticated \
		--platform managed \
		--region $(GACHA_CLOUD_RUN_REGION) $(GACHA_CLOUD_RUN_SERVICE_NAME) \
		--image $(REGISTRY)

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

.PHONY: build test gotest golint govet coverage migrate
