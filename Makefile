BASE_BRANCH ?= master
BENCH_TESTS ?= .
BENCH_DIR ?= bench_results
WORKTREE_DIR ?= .worktree
CURRENT_BRANCH := $(shell git branch --show-current)
ifeq ($(CURRENT_BRANCH),)
CURRENT_BRANCH := $(shell git rev-parse --abbrev-ref HEAD)
ifeq ($(CURRENT_BRANCH),HEAD)
CURRENT_BRANCH := $(shell git rev-parse --short HEAD)
endif
endif
SAFE_BRANCH := $(shell echo "$(CURRENT_BRANCH)" | tr '/' '-')
SAFE_BASE := $(shell echo "$(BASE_BRANCH)" | tr '/' '-')

GOPATH := $(eval GOPATH := $(shell go env GOPATH))$(GOPATH)
GOLINES_MAX_WIDTH ?= 200

ROOTDIR = $(eval ROOTDIR := $(or $(shell git rev-parse --show-toplevel), $(PWD)))$(ROOTDIR)

.PHONY: all clean test bench bench-compare fuzz format lint profile

all: clean test

clean:
	go clean ./...
	rm -rf \
	    $(ROOTDIR)/*.log \
	    $(ROOTDIR)/*.svg \
	    $(ROOTDIR)/pprof \
	    $(ROOTDIR)/bench \
	    $(ROOTDIR)/$(BENCH_DIR) \
	    $(ROOTDIR)/$(WORKTREE_DIR)

test:
	go test -race -v ./...

fuzz:
	go test -fuzz=FuzzMurmur3 -fuzztime=30s ./...

bench:
	go test -count=5 -timeout=30m -run=NONE -bench=$(BENCH_TESTS) -benchmem ./...

$(BENCH_DIR) $(WORKTREE_DIR):
	mkdir -p $@

$(WORKTREE_DIR)/$(BASE_BRANCH): | $(WORKTREE_DIR)
	git worktree remove -f $@ 2>/dev/null || true
	git worktree add $@ $(BASE_BRANCH)

sync-test-files: | $(WORKTREE_DIR)/$(BASE_BRANCH)
	@git diff --name-only $(BASE_BRANCH)...$(CURRENT_BRANCH) | grep -E '_test\.go$$' | while read -r file; do \
		mkdir -p "$(WORKTREE_DIR)/$(BASE_BRANCH)/$$(dirname "$$file")"; \
		if [ -f "$(ROOTDIR)/$$file" ]; then cp "$(ROOTDIR)/$$file" "$(WORKTREE_DIR)/$(BASE_BRANCH)/$$file"; fi; \
	done || true

bench-compare: $(BENCH_DIR) sync-test-files
	@trap 'git worktree remove --force $(WORKTREE_DIR)/$(BASE_BRANCH) 2>/dev/null || true' EXIT; \
	if [ -z "$(CURRENT_BRANCH)" ] || [ "$(CURRENT_BRANCH)" = "$(BASE_BRANCH)" ]; then \
		echo "Must be on a branch other than $(BASE_BRANCH) to compare." && exit 1; \
	fi; \
	echo "Comparing benchmarks: $(BASE_BRANCH) vs $(CURRENT_BRANCH)"; \
	echo "Running benchmarks on $(CURRENT_BRANCH)..."; \
	go test -count=5 -timeout=30m -run=NONE -bench=$(BENCH_TESTS) -benchmem ./... | tee $(ROOTDIR)/$(BENCH_DIR)/$(SAFE_BRANCH).log; \
	echo "Running benchmarks on $(BASE_BRANCH)..."; \
	go test -C $(WORKTREE_DIR)/$(BASE_BRANCH) -count=5 -timeout=30m -run=NONE -bench=$(BENCH_TESTS) -benchmem ./... | tee $(ROOTDIR)/$(BENCH_DIR)/$(SAFE_BASE).log; \
	echo "Comparing results..."; \
	command -v benchstat > /dev/null || go install golang.org/x/perf/cmd/benchstat@latest; \
	$(GOPATH)/bin/benchstat $(ROOTDIR)/$(BENCH_DIR)/$(SAFE_BASE).log $(ROOTDIR)/$(BENCH_DIR)/$(SAFE_BRANCH).log | tee $(ROOTDIR)/$(BENCH_DIR)/benchstat-$(SAFE_BASE)-$(SAFE_BRANCH)

profile:
	mkdir -p bench pprof
	go test -count=3 -timeout=30m -run=NONE -bench=. -benchmem \
		-o pprof/murmur3-test.bin \
		-cpuprofile pprof/cpu-murmur3.out \
		-memprofile pprof/mem-murmur3.out
	go tool pprof --svg pprof/murmur3-test.bin pprof/cpu-murmur3.out > cpu-murmur3.svg
	go tool pprof --svg pprof/murmur3-test.bin pprof/mem-murmur3.out > mem-murmur3.svg
	mv ./*.svg bench/

format:
	find ./ -type d -name .git -prune -o -type f -regex '.*[^\.pb]\.go' -print | xargs $(GOPATH)/bin/golines -w -m $(GOLINES_MAX_WIDTH)
	find ./ -type d -name .git -prune -o -type f -regex '.*[^\.pb]\.go' -print | xargs $(GOPATH)/bin/gofumpt -w
	find ./ -type d -name .git -prune -o -type f -regex '.*\.go' -print | xargs $(GOPATH)/bin/goimports -w
	go fix ./...

lint:
	golangci-lint run --fix
