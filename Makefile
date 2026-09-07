BINARY_NAME=tuiscale
BUILD_DIR=bin
VERSION ?= 1.0.0
COMMIT ?= $(shell git rev-parse --short HEAD 2>/dev/null || echo "none")
DATE ?= $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")

LDFLAGS=-s -w -X main.version=$(VERSION) -X main.commit=$(COMMIT) -X main.date=$(DATE)

.PHONY: all build build-static test clean install run

all: test build

build:
	@mkdir -p $(BUILD_DIR)
	go build -ldflags="$(LDFLAGS)" -o $(BUILD_DIR)/$(BINARY_NAME) ./cmd/tuiscale
	@cp $(BUILD_DIR)/$(BINARY_NAME) ./$(BINARY_NAME)
	@echo "✓ Binário construído com sucesso em ./$(BINARY_NAME)"

# Compilação estática universal (compatível com qualquer distribuição Linux: glibc, musl, alpine, void, arch, etc.)
build-static:
	@mkdir -p $(BUILD_DIR)
	CGO_ENABLED=0 GOOS=linux go build -ldflags="$(LDFLAGS) -extldflags '-static'" -o $(BUILD_DIR)/$(BINARY_NAME)-linux-amd64 ./cmd/tuiscale
	@cp $(BUILD_DIR)/$(BINARY_NAME)-linux-amd64 ./$(BINARY_NAME)
	@echo "✓ Binário estático universal compilado em ./$(BINARY_NAME)"

test:
	go test -v -race ./...

clean:
	rm -rf $(BUILD_DIR) $(BINARY_NAME)

install: build
	@install -d /usr/local/bin
	@install -m 755 ./$(BINARY_NAME) /usr/local/bin/$(BINARY_NAME)
	@echo "✓ TUIScale instalado em /usr/local/bin/$(BINARY_NAME)"

run: build
	./$(BINARY_NAME)
