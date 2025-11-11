BINDIR := bin
BINNAME := fracta
BUILD_TAGS :=
RUN_ARGS :=

.PHONY: all build run clean runtests generate

all: build

build: generate
	@echo "Building '$(BINNAME)'"
	@mkdir -p $(BINDIR)
	go build -tags="$(BUILD_TAGS)" -o $(BINDIR)/$(BINNAME)

run: generate
	@echo "Running"
	@echo "------------------------------------------\n"
	@go run -tags="$(BUILD_TAGS)" ./ $(RUN_ARGS)

clean:
	@echo "Cleaning up"
	rm -rf $(BINDIR)
	rm -f $(ERRORFILE_OUT)

runtests: generate
	@echo "Running tests"
	@go test ./test


# Code generation targets below this point

ERRORFILE_OUT := internal/diag/errors.go
ERRORFILE_CONFIG := gen/error_config.yaml
ERRORFILE_GEN := gen/gen_errors.go

# Generate code if YAML changed or output missing
$(ERRORFILE_OUT): $(ERRORFILE_CONFIG) $(ERRORFILE_GEN)
	@echo "[GEN] Generating error definitions..."
	go generate ./...


TOKENTYPE_STRING := internal/token/tokentype_string.go
TOKENTYPE_SRC := internal/token/token.go

# Generate TokenType strings
$(TOKENTYPE_STRING): $(TOKENTYPE_SRC)
	@echo "[GEN] Generating TokenType string definitions..."
	stringer -type TokenType ./internal/token

generate: $(ERRORFILE_OUT) $(TOKENTYPE_STRING)
