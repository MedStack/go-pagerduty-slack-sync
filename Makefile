
TAG ?= $$(git describe --tags)

build:
	@env GOMODULE111=on find ./cmd/* -maxdepth 1 -type d -exec go build "{}" \;

install-lint:
	@go get -u golang.org/x/lint/golint

install-deps: install-lint

lint:
	@golint ./...

vet:
	@go vet -v ./...

check: lint vet

test:
	@go test -v ./...

test-cover:
	@go test -cover -coverprofile ./coverage.out -v ./...

docker-build:
	@docker build -t medstackcibaa9.azurecr.io/pagerduty-slack-sync -f build/package/Dockerfile .

docker-publish:
	@docker login

ci: install-deps build check test

.PHONY: build install-deps install-lint lint vet check test