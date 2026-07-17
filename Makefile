PKG := github.com/pfolta/cdrdao2audio

# Build metadata injected into the binary via -ldflags.
BUILD_DATE := $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")
VERSION := $(shell git describe --match "v[0-9]*" --dirty="-m" --always --tags || echo "dev")

BUILD_DIR ?= ./build

# Use the requested target platform if provided.
# Otherwise, default to the platform of the build host.
GOOS ?= $(shell go env GOOS)
GOARCH ?= $(shell go env GOARCH)

ifeq ($(GOOS),windows)
	BINARY := cdrdao2audio.exe
else
	BINARY := cdrdao2audio
endif

GO_LDFLAGS := \
	-X "$(PKG).buildDate=$(BUILD_DATE)" \
	-X "$(PKG).version=$(VERSION)"

.PHONY: default
default: release

.PHONY: build
build:
	mkdir -p "$(BUILD_DIR)/bin"

	CGO_ENABLED=0 \
	GOOS="$(GOOS)" \
	GOARCH="$(GOARCH)" \
	go build \
		-ldflags "$(GO_LDFLAGS)" \
		-o "$(BUILD_DIR)/bin/$(BINARY)" \
		./cmd/cdrdao2audio

.PHONY: clean
clean:
	rm -rf "$(BUILD_DIR)"

.PHONY: install
install:
	CGO_ENABLED=0 \
	go install \
		-ldflags "$(GO_LDFLAGS)" \
		./cmd/cdrdao2audio

.PHONY: lint
lint:
	gofmt -d .
	go vet ./...

.PHONY: release
release: lint test build

.PHONY: test
test:
	mkdir -p "$(BUILD_DIR)/tests"
	mkdir -p "$(BUILD_DIR)/reports"

	go test \
		-race \
		-covermode=atomic \
		-coverprofile="$(BUILD_DIR)/tests/coverage.out" \
		-v ./...

	go tool \
		cover \
		-html="$(BUILD_DIR)/tests/coverage.out" \
		-o "$(BUILD_DIR)/reports/coverage.html"

.PHONY: version
version:
	@echo "$(VERSION)"
