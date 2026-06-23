
VERSION := $(shell git describe --match "v[0-9]*" --dirty="-m" --always --tags || echo "dev")
BUILD_DATE := $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")

BUILD_DIR := ./build

GO_LDFLAGS := \
	-X "main.Version=$(VERSION)" \
	-X "main.BuildDate=$(BUILD_DATE)"

ifeq ($(OS),Windows_NT)
	BINARY := cdrdao2wav.exe
	CLEAN_CMD := rmdir /s /q
else
	BINARY := cdrdao2wav
	CLEAN_CMD := rm -rf
endif

.PHONY: build
build:
	go build -ldflags "$(GO_LDFLAGS)" -o "$(BUILD_DIR)/bin/$(BINARY)" .

.PHONY: clean
clean:
	$(CLEAN_CMD) $(BUILD_DIR)

.PHONY: lint
lint:
	go vet ./...

.PHONY: release
release: lint test build

.PHONY: test
test:
	mkdir -p $(BUILD_DIR)/tests
	mkdir -p $(BUILD_DIR)/reports
	go test -race -v ./... -coverprofile=$(BUILD_DIR)/tests/coverage.out
	go tool cover -html=$(BUILD_DIR)/tests/coverage.out -o $(BUILD_DIR)/reports/coverage.html
