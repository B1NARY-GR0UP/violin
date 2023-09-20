GO := go

ifeq ($(OS), Windows_NT)
	RM := del
else
	RM := rm -f
endif

SRC_FILES := $(wildcard *.go)
TEST_FILES := $(wildcard *_test.go)

TEST_CMD := $(GO) test -v ./...

COVERAGE_CMD := $(GO) test -coverprofile="coverage.out" ./...

BENCHMARK_CMD := $(GO) test -bench=. ./...

default: test

test:
	@$(TEST_CMD)

coverage:
	@$(COVERAGE_CMD)
	@$(GO) tool cover -html="coverage.out" -o "coverage.html"

benchmark:
	@$(BENCHMARK_CMD)

clean:
	@$(RM) coverage.out coverage.html

format:
	@gofumpt -e -d -w -extra .

.PHONY: test coverage benchmark clean format