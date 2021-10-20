GO ?= go
GOPATH ?= $(HOME)/go
GO_INSTALL_OPTS ?=
GO_TEST_OPTS ?= -test.timeout=30s
GOMOD_DIRS ?= $(sort $(call novendor,$(dir $(call rwildcard,*,*/go.mod go.mod))))
GOCOVERAGE_FILE ?= ./coverage.txt
GOTESTJSON_FILE ?= ./go-test.json
GOBUILDLOG_FILE ?= ./go-build.log
GOINSTALLLOG_FILE ?= ./go-install.log

.PHONY: unittest
unittest:
	@echo "mode: atomic" > /tmp/gocoverage
		@rm -f $(GOTESTJSON_FILE)
		@set -e; for dir in $(GOMOD_DIRS); do (set -e; (set -euf pipefail; \
			cd $$dir; \
			(($(GO) test ./... $(GO_TEST_OPTS) -cover -coverprofile=/tmp/profile.out -covermode=atomic -race -json && touch $@.ok) | tee -a $(GOTESTJSON_FILE) 3>&1 1>&2 2>&3 | tee -a $(GOBUILDLOG_FILE); \
		  ); \
		  rm $@.ok 2>/dev/null || exit 1; \
		  if [ -f /tmp/profile.out ]; then \
			cat /tmp/profile.out | sed "/mode: atomic/d" >> /tmp/gocoverage; \
			rm -f /tmp/profile.out; \
		  fi)); done
		@mv /tmp/gocoverage $(GOCOVERAGE_FILE)