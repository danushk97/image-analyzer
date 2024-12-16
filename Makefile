# Dir where build binaries are generated. The dir should be gitignored
BUILD_OUT_DIR := "bin/"

API_OUT       := "bin/api"
API_MAIN_FILE := "cmd/server/main.go"

MIGRATION_OUT       := "bin/migration"
MIGRATION_MAIN_FILE := "cmd/migration/main.go"

# go binary. Change this to experiment with different versions of go.
GO       = go

MODULE   = $(shell $(GO) list -m)
SERVICE  = $(shell basename $(MODULE))

# Proto gen info
PROTO_ROOT := proto/
RPC_ROOT := rpc/

# Fetch OS info
GOVERSION=$(shell go version)
UNAME_OS=$(shell go env GOOS)
UNAME_ARCH=$(shell go env GOARCH)

# This is the only variable that ever should change.
# This can be a branch, tag, or commit.
BUF_VERSION := v1.5.0
PROTOC_GEN_GO_VERSION := v1.3.2
PROTOC_GEN_TWIRP_VERSION := v5.10.1

.PHONY: all
all: build

.PHONY: build-info
build-info:
	@echo "\nBuild Info:\n"
	@echo "\t\033[33mOS\033[0m: $(UNAME_OS)"
	@echo "\t\033[33mArch\033[0m: $(UNAME_ARCH)"
	@echo "\t\033[33mGo Version\033[0m: $(GOVERSION)\n"

.PHONY: go-build-api ## Build the binary file for API server
go-build-api:
	@CGO_ENABLED=0 GOOS=$(UNAME_OS) GOARCH=$(UNAME_ARCH) go build -v -o $(API_OUT) $(API_MAIN_FILE)

.PHONY: go-build-migration ## Build the binary file for database migrations
go-build-migration:
	@CGO_ENABLED=0 GOOS=$(UNAME_OS) GOARCH=$(UNAME_ARCH) go build -v -o $(MIGRATION_OUT) $(MIGRATION_MAIN_FILE)

