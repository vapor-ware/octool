#
# octool : Basic tool to check OpenConfig connectivity
#

BIN_NAME    := octool
BIN_VERSION := v0.1.0-rc0

LDFLAGS := -w -s

.PHONY: build build-linux clean fmt github-tag lint version help

.DEFAULT_GOAL := help


build:  ## Build the binary
	go build -ldflags "${LDFLAGS}" -o ${BIN_NAME}

build-linux:  ## Build the binary for linux amd64
	GOOS=linux GOARCH=amd64 go build -ldflags "${LDFLAGS}" -o ${BIN_NAME} .

clean:  ## Remove temporary files
	go clean -v
	rm -rf dist


fmt:  ## Run goimports on all go files
	find . -name '*.go' -not -wholename './vendor/*' -not -wholename '*.pb.go' | while read -r file; do goimports -w "$$file"; done

github-tag:  ## Create and push a tag with the current version
	git tag -a ${BIN_VERSION} -m "${BIN_NAME} version ${BIN_VERSION}"
	git push -u origin ${BIN_VERSION}

lint:  ## Lint project source files
	golint -set_exit_status ./pkg/...


version:  ## Print the version of the binary
	@echo "${BIN_VERSION}"

help:  ## Print usage information
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "\033[36m%-15s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST) | sort
