GOCMD=go
GOBUILD=${GOCMD} build
GOCLEAN=${GOCMD} clean
GOGET=${GOCMD} get
GOTEST=${GOCMD} test
GORUN=${GOCMD} run 
BINARY_NAME=guest-agent-test-extension.exe

.PHONY: all
all: clean build

.PHONY: test
test: 
	${GOTEST} -v


.PHONY: build
build: deps
	${GOBUILD} -o ${BINARY_NAME} -v

.PHONY: deps
deps:
	${GOGET} -u "github.com/go-kit/kit/log"
	${GOGET} -u "github.com/Azure/azure-extension-foundation/sequence"
	${GOGET} -u "github.com/Azure/azure-extension-foundation/settings"
	${GOGET} -u "github.com/Azure/azure-extension-foundation/status"

.PHONY: clean
clean:
	${GOCLEAN}


.PHONY: help
help:
	@echo "TODO"
