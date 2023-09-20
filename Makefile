GO := go

SRC_FILES := $(wildcard *.go)
TEST_FILES := $(wildcard *_test.go)
COVERAGE_FILE := coverage.out

TEST_CMD := $(GO) test -v ./...

COVERAGE_CMD := $(GO) test -coverprofile=$(COVERAGE_FILE) ./...

BENCHMARK_CMD := $(GO) test -bench=. ./...

default: test

test:
	@$(TEST_CMD)

coverage:
	@$(COVERAGE_CMD)
	@$(GO) tool cover -html=$(COVERAGE_FILE)

benchmark:
	@$(BENCHMARK_CMD)

.PHONY: test coverage benchmark