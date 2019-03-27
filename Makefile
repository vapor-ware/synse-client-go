#
# Golang Client for Synse Server HTTP API
#

PKG_VERSION := 0.0.1

HAS_LINT := $(shell which gometalinter)
HAS_DEP  := $(shell which dep)

#
# Development Targets
#

.PHONY: clean
clean:  ## Remove temporary files
	go clean -v

.PHONY: cover
cover: test  ## Run unit tests and open the coverage report
	go tool cover -html=coverage.out

.PHONY: dep
dep:  ## Ensure and prune dependencies
ifndef HAS_DEP
	go get -u github.com/golang/dep/cmd/dep
endif
	dep ensure -v

.PHONY: fmt
fmt:  ## Run goimports on all go source files
	find . -name '*.go' -not -wholename './vendor/*' | while read -r file; do goimports -w "$$file"; done

.PHONY: github-tag
github-tag:  ## Create and push a tag with the current client version
	git tag -a ${PKG_VERSION} -m "Synse Go Client version ${PKG_VERSION}"
	git push -u origin ${PKG_VERSION}

.PHONY: lint
lint:  ## Lint project source files
	gometalinter ./... \
		--tests \
		--vendor \
		--sort=path --sort=line \
		--aggregate \
		--deadline=5m \
		-e $$(go env GOROOT)

.PHONY: setup
setup:  ## Install the build and development dependencies and set up vendoring
	go get -u github.com/alecthomas/gometalinter
	go get -u github.com/golang/dep/cmd/dep
	gometalinter --install
ifeq (,$(wildcard ./Gopkg.toml))
	dep init
endif
	@$(MAKE) dep

.PHONY: test
test:  ## Run all unit tests
	@ # Note: this requires go1.10+ in order to do multi-package coverage reports
	go test -race -coverprofile=coverage.out -covermode=atomic ./...

.PHONY: version
version:  ## Print the version of the client
	@echo "$(PKG_VERSION)"

.PHONY: help
help:  ## Print usage information
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "\033[36m%-15s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST) | sort

.DEFAULT_GOAL := help


#
# CI Targets
#

.PHONY: ci-check-version
ci-check-version:
	PKG_VERSION=$(PKG_VERSION) ./bin/ci/check_version.sh
