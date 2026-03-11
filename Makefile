.PHONY: build test test-unit test-race test-e2e lint format coverage doc-design dup complexity release sec help

GO := go
BUN := bun
GOLINT := golangci-lint
ROOT_DIR := $(PWD)
TMP_DIR := $(ROOT_DIR)/tmp
GO_CACHE_DIR := $(ROOT_DIR)/.gocache
GO_MOD_CACHE_DIR := $(ROOT_DIR)/.gomodcache
GO_LINT_CACHE_DIR := $(ROOT_DIR)/.golangci-lint-cache
E2E_BIN_DIR := $(ROOT_DIR)/.e2e-bin
GO_PACKAGES := ./...
GO_ENV := GOTOOLCHAIN=local GOCACHE=$(GO_CACHE_DIR) GOMODCACHE=$(GO_MOD_CACHE_DIR)
BUN_ENV := TMPDIR=$(TMP_DIR)
BIOME := $(BUN_ENV) $(BUN) run biome
GOLINT_ENV := $(GO_ENV) GOLANGCI_LINT_CACHE=$(GO_LINT_CACHE_DIR)
COVER_PROFILE := $(TMP_DIR)/test-unit.coverage.out
COVER_HTML := $(TMP_DIR)/test-unit.coverage.html
VERSION ?= dev
COMMIT ?= unknown
BUILD_DATE ?= unknown
GO_LDFLAGS := -X github.com/flarebyte/baldrick-seer/internal/buildinfo.Version=$(VERSION) -X github.com/flarebyte/baldrick-seer/internal/buildinfo.Commit=$(COMMIT) -X github.com/flarebyte/baldrick-seer/internal/buildinfo.Date=$(BUILD_DATE)

build:
	mkdir -p $(TMP_DIR)
	mkdir -p $(E2E_BIN_DIR)
	$(GO_ENV) $(GO) build -ldflags "$(GO_LDFLAGS)" -o $(E2E_BIN_DIR)/seer ./cmd/seer

test: test-unit test-e2e

test-unit:
	mkdir -p $(TMP_DIR)
	$(GO_ENV) $(GO) test -v -coverprofile=$(COVER_PROFILE) -covermode=count $(GO_PACKAGES)
	$(GO_ENV) $(GO) tool cover -func=$(COVER_PROFILE)

test-race:
	mkdir -p $(TMP_DIR)
	$(GO_ENV) $(GO) test -race $(GO_PACKAGES)

coverage: test-unit
	$(GO_ENV) $(GO) tool cover -html=$(COVER_PROFILE) -o $(COVER_HTML)
	@printf "Coverage HTML: %s\n" "$(COVER_HTML)"

test-e2e: build
	mkdir -p $(TMP_DIR)
	$(BUN_ENV) $(BUN) install
	$(BUN_ENV) $(BUN) test ./e2e

lint:
	mkdir -p $(TMP_DIR)
	$(BUN_ENV) $(BUN) install
	$(BIOME) check .
	$(GO_ENV) $(GO) vet $(GO_PACKAGES)
	$(GOLINT_ENV) $(GOLINT) run

format:
	mkdir -p $(TMP_DIR)
	$(BUN_ENV) $(BUN) install
	find . -type f -name '*.go' \
		-not -path './.git/*' \
		-not -path './.gocache/*' \
		-not -path './.gomodcache/*' \
		-not -path './.e2e-bin/*' \
		-not -path './node_modules/*' \
		-print0 | xargs -0 -r gofmt -w
	$(BIOME) format --write .

doc-design:
	mkdir -p doc/design
	flyb validate --config doc/design-meta/app.cue
	flyb generate markdown --config doc/design-meta/app.cue
	flyb validate --config doc/design-meta/flows.cue
	flyb generate markdown --config doc/design-meta/flows.cue

dup:
	npx jscpd --format go --min-lines 10 --ignore "**/.gomodcache/**,**/.gocache/**,**/.e2e-bin/**,**/node_modules/**,**/dist/**" --gitignore .
	npx jscpd --format typescript --min-lines 10 --gitignore .

complexity:
	scc --sort complexity --by-file -i go . | head -n 15
	scc --sort complexity --by-file -i ts . | head -n 15

release:
	$(BUN_ENV) $(BUN) run release-go.ts

sec:
	semgrep scan --config auto

help:
	@printf "Targets:\n"
	@printf "  build       Build the seer binary into .e2e-bin/ with deterministic ldflags.\n"
	@printf "  test        Run unit and E2E tests.\n"
	@printf "  test-unit   Run verbose Go unit tests and print coverage.\n"
	@printf "  test-race   Run Go tests with the race detector.\n"
	@printf "  test-e2e    Build the CLI and run Bun E2E smoke tests.\n"
	@printf "  coverage    Generate the unit-test coverage summary and HTML report.\n"
	@printf "  lint        Run Biome checks for E2E/tooling files and go vet.\n"
	@printf "  format      Format Go files and Biome-managed files.\n"
	@printf "  doc-design  Regenerate design docs from flyb configs.\n"
	@printf "  dup         Run duplicate code detection.\n"
	@printf "  complexity  Show top Go and TypeScript files by complexity.\n"
	@printf "  release     Run the release helper script.\n"
	@printf "  sec         Run Semgrep security scan.\n"
	@printf "  help        Show this help message.\n"
