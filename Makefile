#
# Golang Client for Synse Server HTTP API
#

VERSION := 0.0.1

.PHONY: clean
clean:  ## Remove temporary files
	go clean -v

.PHONY: fmt
fmt:  ## Run goimports on all go source files
	find . -name '*.go' -not -wholename './vendor/*' | while read -r file; do goimports -w "$$file"; done

.PHONY: github-tag
github-tag:  ## Create and push a tag with the current client version
	git tag -a ${VERSION} -m "Synse HTTP Go Client version ${VERSION}"
	git push -u origin ${VERSION}

.PHONY: lint
lint:  ## Lint project source files
	gometalinter ./... \
		--tests \
		--vendor \
		--sort=path --sort=line \
		--aggregate \
		--deadline=5m \
		-e $$(go env GOROOT)

.PHONY: version
version:  ## Print the version of the client
	@echo "$(VERSION)"

.PHONY: help
help:  ## Print usage information
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "\033[36m%-15s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST) | sort

.DEFAULT_GOAL := help
