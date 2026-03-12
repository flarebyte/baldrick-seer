.PHONY: build build-go build-dist test test-go test-unit test-race test-e2e lint lint-go lint-e2e format format-go format-e2e coverage coverage-go doc-design dup complexity release sec help

GO := go
BUN := bun
GOLINT := golangci-lint
ROOT_DIR := $(PWD)
TMP_DIR := $(ROOT_DIR)/tmp
GO_CACHE_DIR := $(ROOT_DIR)/.gocache
GO_MOD_CACHE_DIR := $(ROOT_DIR)/.gomodcache
GO_LINT_CACHE_DIR := $(ROOT_DIR)/.golangci-lint-cache
E2E_BIN_DIR := $(ROOT_DIR)/.e2e-bin
BUILD_DIR := $(ROOT_DIR)/build
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

build: build-go build-dist

build-go:
	mkdir -p $(TMP_DIR)
	mkdir -p $(E2E_BIN_DIR)
	$(GO_ENV) $(GO) build -ldflags "$(GO_LDFLAGS)" -o $(E2E_BIN_DIR)/seer ./cmd/seer

build-dist:
	mkdir -p $(TMP_DIR)
	mkdir -p $(BUILD_DIR)
	$(BUN_ENV) $(BUN) run build-go.ts

test: test-go test-e2e

test-go: test-unit

test-unit:
	mkdir -p $(TMP_DIR)
	$(GO_ENV) $(GO) test -v -coverprofile=$(COVER_PROFILE) -covermode=count $(GO_PACKAGES)
	$(GO_ENV) $(GO) tool cover -func=$(COVER_PROFILE)

test-race:
	mkdir -p $(TMP_DIR)
	$(GO_ENV) $(GO) test -race $(GO_PACKAGES)

coverage: coverage-go

coverage-go: test-unit
	$(GO_ENV) $(GO) tool cover -html=$(COVER_PROFILE) -o $(COVER_HTML)
	@printf "Coverage HTML: %s\n" "$(COVER_HTML)"

test-e2e: build-go
	mkdir -p $(TMP_DIR)
	$(BUN_ENV) $(BUN) install
	$(BUN_ENV) $(BUN) test ./e2e

lint: lint-go lint-e2e

lint-go:
	mkdir -p $(TMP_DIR)
	$(GO_ENV) $(GO) vet $(GO_PACKAGES)
	$(GOLINT_ENV) $(GOLINT) run

lint-e2e:
	mkdir -p $(TMP_DIR)
	$(BUN_ENV) $(BUN) install
	$(BIOME) check .

format: format-go format-e2e

format-go:
	mkdir -p $(TMP_DIR)
	find . -type f -name '*.go' \
		-not -path './.git/*' \
		-not -path './.gocache/*' \
		-not -path './.gomodcache/*' \
		-not -path './.e2e-bin/*' \
		-not -path './node_modules/*' \
		-print0 | xargs -0 -r gofmt -w

format-e2e:
	mkdir -p $(TMP_DIR)
	$(BUN_ENV) $(BUN) install
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
	@printf "  build        Build the E2E binary and release artifacts.\n"
	@printf "  build-go     Build the Go CLI into .e2e-bin/ for local and E2E use.\n"
	@printf "  build-dist   Build multi-platform release binaries into build/.\n"
	@printf "  test         Run Go tests and Bun E2E tests.\n"
	@printf "  test-go      Run Go test targets.\n"
	@printf "  test-unit    Run verbose Go tests and print coverage summary.\n"
	@printf "  test-race    Run Go tests with the race detector.\n"
	@printf "  test-e2e     Build the CLI and run Bun E2E tests.\n"
	@printf "  coverage     Generate the coverage HTML report.\n"
	@printf "  coverage-go  Generate the Go coverage HTML report from test-unit output.\n"
	@printf "  lint         Run Go linting and Biome checks.\n"
	@printf "  lint-go      Run go vet and golangci-lint.\n"
	@printf "  lint-e2e     Run Biome checks for TypeScript and tooling files.\n"
	@printf "  format       Format Go and Biome-managed files.\n"
	@printf "  format-go    Format Go files with gofmt.\n"
	@printf "  format-e2e   Format TypeScript and tooling files with Biome.\n"
	@printf "  doc-design   Regenerate design docs from flyb configs.\n"
	@printf "  dup          Run duplicate code detection.\n"
	@printf "  complexity   Show top Go and TypeScript files by complexity.\n"
	@printf "  release      Run the local release helper script.\n"
	@printf "  sec          Run Semgrep security scan.\n"
	@printf "  help         Show this help message.\n"
