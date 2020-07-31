SHELL := /bin/bash

GOCMD=go
GOBUILD=${GOCMD} build
GOCLEAN=${GOCMD} clean
GOGET=${GOCMD} get
GOTEST=${GOCMD} test
GORUN=${GOCMD} run 

BINARY_NAME=bin/guest-agent-test-extension
WINDOWS_BIN=$(BINARY_NAME)_windows.exe
LINUX_BIN=$(BINARY_NAME)_linux

.PHONY: all
all: clean build_all

.PHONY: test
test: 
	${GOTEST} -v


.PHONY: build_all
build_all: build_windows build_linux


.PHONY: build_windows
build_windows: deps
	$(GOCMD) env GOOS=windows 
	${GOBUILD} -o ${WINDOWS_BIN} -v

.PHONY: build_linux
build_linux: deps
	$(GOCMD) env GOOS=linux
	${GOBUILD} -o ${LINUX_BIN} -v



.PHONY: deps
deps:
	${GOGET} -u "github.com/Azure/azure-extension-foundation/sequence"
	${GOGET} -u "github.com/Azure/azure-extension-foundation/settings"
	${GOGET} -u "github.com/Azure/azure-extension-foundation/status"

.PHONY: clean
clean:
	${GOCLEAN}
	

help:
	@echo "TODO"
