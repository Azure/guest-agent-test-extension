GOCMD=go
GOBUILD=${GOCMD} build
GOCLEAN=${GOCMD} clean
GOGET=${GOCMD} get
GOTEST=${GOCMD} test
GORUN=${GOCMD} run 

BINARY_NAME=bin/GuestAgentTestExtension
WINDOWS_BIN=$(BINARY_NAME)_windows.exe
LINUX_BIN=$(BINARY_NAME)_linux
SERVICE_SCRIPT_BIN=bin/gatestext_script_linux

.PHONY: all
all: clean build_all

.PHONY: test
test: 
	${GOTEST} -v


.PHONY: build_all
build_all: build_windows build_linux

.PHONY: build_windows
build_windows:
	$(GOCMD) env -w GOOS=windows 
	${GOBUILD} -o ${WINDOWS_BIN} ./main/

.PHONY: build_linux
build_linux:
	$(GOCMD) env -w GOOS=linux
	${GOBUILD} -o  ${LINUX_BIN} ./main/
	${GOBUILD} -o  ${SERVICE_SCRIPT_BIN} ./services/

.PHONY: build_all_with_deps
build_all_with_deps: deps build_all

.PHONY: build_linux_with_deps
build_linux_with_deps: deps build_linux

.PHONY: build_windows_with_deps
build_windows_with_deps: deps build_windows

.PHONY: deps
deps:
	${GOGET} -u "github.com/Azure/azure-extension-foundation/sequence"
	${GOGET} -u "github.com/Azure/azure-extension-foundation/settings"
	${GOGET} -u "github.com/Azure/azure-extension-foundation/status"
	${GOGET} -u "github.com/pkg/errors"

.PHONY: clean
clean:
	${GOCLEAN}
	

help:
	@echo "TODO"