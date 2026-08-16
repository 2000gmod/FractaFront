# Target platforms (OS/ARCH) - override with PLATFORMS="linux/amd64 linux/arm64"
PLATFORMS ?= linux/amd64

# Convert platforms to safe target names (e.g., linux/amd64 -> build-linux_amd64)
PLATFORM_TARGETS := $(foreach p,$(PLATFORMS),build-$(subst /,_,$p))

# Set to 1 to use Docker buildx, otherwise native go build (host arch only)
DOCKER_BUILD ?= 1

# Output directory (relative to project root)
BUILD ?= build

# Name of the resulting binary
BIN_NAME ?= fracta

# Go build tags (e.g., TAGS="-tags=llvm19,b_llvm")
TAGS ?= -tags=llvm19,b_llvm

# Host platform detection (used to guard native cross-compilation)
HOST_OS     := $(shell go env GOOS)
HOST_ARCH   := $(shell go env GOARCH)
HOST_PLATFORM := $(HOST_OS)/$(HOST_ARCH)

# Buildx builder name (for isolated multi-arch setup)
BUILDX_BUILDER ?= fracta-builder

# BuildKit cache directories (persists across builds)
CACHE_DIR ?= /tmp/buildx-cache

.PHONY: all build clean setup help

all: build

setup:
	@echo "Setting up Docker Buildx builder..."
	@docker buildx create --name $(BUILDX_BUILDER) --use --bootstrap 2>/dev/null || docker buildx use $(BUILDX_BUILDER)
	@echo "Registering QEMU emulators for cross-architecture builds..."
	@docker run --privileged --rm tonistiigi/binfmt --install all
	@echo "Setup complete! You can now run 'make build'."

build: $(PLATFORM_TARGETS)

build-%:
	@set -e; \
	platform_underscore=$*; \
	platform=$$(echo $$platform_underscore | tr '_' '/'); \
	platform_dir=$$(echo $$platform | tr '/' '_'); \
	mkdir -p $(BUILD)/$$platform_dir; \
	output_subdir=$(BUILD)/$$platform_dir; \
	if [ "$(DOCKER_BUILD)" = "1" ]; then \
		echo ">>> Building with Docker for $$platform (using builder: $(BUILDX_BUILDER))"; \
		docker buildx build \
			--builder $(BUILDX_BUILDER) \
			--platform $$platform \
			--build-arg BIN_NAME=$(BIN_NAME) \
			--build-arg TAGS="$(TAGS)" \
			--cache-from type=local,src=$(CACHE_DIR) \
			--cache-to type=local,dest=$(CACHE_DIR),mode=max \
			--output type=local,dest=./$$output_subdir \
			. ; \
		if [ -f "$$output_subdir/$(BIN_NAME)" ]; then \
			mv "$$output_subdir/$(BIN_NAME)" "$$output_subdir/$(BIN_NAME)-$$platform_dir"; \
		fi; \
	else \
		os=$$(echo $$platform | cut -d'/' -f1); \
		arch=$$(echo $$platform | cut -d'/' -f2); \
		echo ">>> Building natively for GOOS=$$os GOARCH=$$arch"; \
		if [ "$$platform" != "$(HOST_PLATFORM)" ]; then \
			echo "ERROR: Native cross-compilation with CGO/LLVM is not supported."; \
			echo "Set DOCKER_BUILD=1 or build only for your host platform ($(HOST_PLATFORM))."; \
			exit 1; \
		fi; \
		GOOS=$$os GOARCH=$$arch go build $(TAGS) -ldflags="$(LDFLAGS_FLAGS)" -o ./$$output_subdir/$(BIN_NAME)-$$platform_dir . ; \
	fi; \
	echo "Built: ./$$output_subdir/$(BIN_NAME)-$$platform_dir"

clean:
	@rm -rf $(BUILD)
	@echo "Removed $(BUILD)/"

help:
	@echo "Available targets:"
	@echo "  setup         - Initialize Buildx builder and QEMU emulators (run once)"
	@echo "  build         - Build for all platforms (parallel). Use 'make -j2 build'"
	@echo "  clean         - Remove the $(BUILD) directory"
	@echo ""
	@echo "Variables (override via make VAR=value):"
	@echo "  PLATFORMS     = $(PLATFORMS)"
	@echo "  DOCKER_BUILD  = $(DOCKER_BUILD)  (1 = Docker, 0 = native go)"
	@echo "  BUILD         = $(BUILD)"
	@echo "  BIN_NAME      = $(BIN_NAME)"
	@echo "  TAGS          = $(TAGS)"
	@echo "  VERSION       = $(VERSION)"
	@echo "  CACHE_DIR     = $(CACHE_DIR)"
