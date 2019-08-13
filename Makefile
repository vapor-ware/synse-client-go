#
# synse-client-go # 
PKG_VERSION := 0.0.1


.PHONY: build
build:  ## Build the library - this will verify correctness but will not produce a binary
	CGO_ENABLED=0 go build ./...

.PHONY: clean
clean:  ## Remove temporary files and build artifacts
	go clean -v
	rm -rf dist
	rm -f ${BIN_NAME} coverage.out

.PHONY: cover
cover: test  ## Run unit tests and open the coverage report
	go tool cover -html=coverage.out

.PHONY: fmt
fmt:  ## Run goimports on all go source files
	find . -name '*.go' -not -wholename './vendor/*' | while read -r file; do goimports -w "$$file"; done

.PHONY: github-tag
github-tag:  ## Create and push a tag with the current client version
	git tag -a ${PKG_VERSION} -m "Synse Go Client v${PKG_VERSION}"
	git push -u origin ${PKG_VERSION}

.PHONY: lint
lint:  ## Lint project source files
	golint -set_exit_status ./...

.PHONY: test
test:  ## Run unit tests
	@ # Note: this requires go1.10+ in order to do multi-package coverage reports
	go test -short -race -coverprofile=coverage.out -covermode=atomic ./...

.PHONY: test-integration
test-integration:  ## Run integration tests
	# FIXME - need to clean up stale containers using `docker rm -f $(docker ps -aq)`
	# every time before running the tests so they won't fail
	docker-compose -f compose/server.yml up -d
	# have to wait at least 30 seconds for the emulated health checks to be fully populated
	sleep 30
	go test -race -cover -run Integration ./... || (docker-compose -f compose/server.yml stop; exit 1)
	docker-compose -f compose/server.yml down

.PHONY: version
version:  ## Print the version of the client
	@echo "$(PKG_VERSION)"

.PHONY: help
help:  ## Print usage information
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "\033[36m%-16s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST) | sort

.DEFAULT_GOAL := help
