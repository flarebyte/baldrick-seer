.PHONY: build test test-unit test-e2e lint format doc-design dup complexity release sec help

GO := go
BUN := bun
GO_ENV := GOTOOLCHAIN=local GOCACHE=$(PWD)/.gocache GOMODCACHE=$(PWD)/.gomodcache
BUN_ENV := TMPDIR=$(PWD)/tmp
BIOME := $(BUN_ENV) $(BUN) run biome

build:
	mkdir -p tmp
	mkdir -p .e2e-bin
	$(GO_ENV) $(GO) build -o .e2e-bin/seer ./cmd/seer

test: test-unit test-e2e

test-unit:
	$(GO_ENV) $(GO) test ./...

test-e2e: build
	mkdir -p tmp
	$(BUN_ENV) $(BUN) install
	$(BUN_ENV) $(BUN) test ./e2e

lint:
	mkdir -p tmp
	$(BUN_ENV) $(BUN) install
	$(BIOME) check .
	$(GO_ENV) $(GO) vet ./...

format:
	mkdir -p tmp
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
	@printf "  build       Build the seer binary into .e2e-bin/.\n"
	@printf "  test        Run unit and E2E tests.\n"
	@printf "  test-unit   Run Go unit tests.\n"
	@printf "  test-e2e    Build the CLI and run Bun E2E smoke tests.\n"
	@printf "  lint        Run Biome checks for E2E/tooling files and go vet.\n"
	@printf "  format      Format Go files and Biome-managed files.\n"
	@printf "  doc-design  Regenerate design docs from flyb configs.\n"
	@printf "  dup         Run duplicate code detection.\n"
	@printf "  complexity  Show top Go and TypeScript files by complexity.\n"
	@printf "  release     Run the release helper script.\n"
	@printf "  sec         Run Semgrep security scan.\n"
	@printf "  help        Show this help message.\n"
