#
# synse-client-go
#

PKG_VERSION := v1.0.0


.PHONY: build
build:  ## Build the library - this will verify correctness but will not produce a binary
	CGO_ENABLED=0 go build ./...

.PHONY: clean
clean:  ## Remove temporary files and build artifacts
	go clean -v
	rm -rf dist coverage.out

.PHONY: cover
cover: test  ## Run unit tests and open the coverage report
	go tool cover -html=coverage.out

.PHONY: fmt
fmt:  ## Run goimports on all go source files
	find . -name '*.go' -not -wholename './vendor/*' | while read -r file; do goimports -w "$$file"; done

.PHONY: github-tag
github-tag:  ## Create and push a tag with the current client version
	git tag -a ${PKG_VERSION} -m "Synse Go Client version ${PKG_VERSION}"
	git push -u origin ${PKG_VERSION}

.PHONY: lint
lint:  ## Lint project source files
	golint -set_exit_status ./...

.PHONY: test
test: test-unit test-integration  ## Run all tests (unit, integration)

.PHONY: test-unit
test-unit:  ## Run unit tests
	go test -short -race -coverprofile=coverage.out -covermode=atomic ./...

.PHONY: test-integration
test-integration: test-integration-http test-integration-websocket  ## Run integration tests

.PHONY: test-integration-http
test-integration-http:  ## Run integration tests for http client
	-docker-compose -f compose/server.yml rm -fsv
	docker-compose -f compose/server.yml up -d
	sleep 6
	go test -race -cover -run IntegrationHTTP ./... || (docker-compose -f compose/server.yml stop; exit 1)
	docker-compose -f compose/server.yml down

.PHONY: test-integration-websocket
test-integration-websocket:  ## Run integration tests for websocket client
	-docker-compose -f compose/server.yml rm -fsv
	docker-compose -f compose/server.yml up -d
	sleep 6
	go test -race -cover -run IntegrationWebSocket ./... || (docker-compose -f compose/server.yml stop; exit 1)
	docker-compose -f compose/server.yml down

.PHONY: version
version:  ## Print the package version
	@echo "$(PKG_VERSION)"

.PHONY: help
help:  ## Print usage information
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "\033[36m%-26s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST) | sort

.DEFAULT_GOAL := help
