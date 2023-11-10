GO_SRCS := $(shell find . -type f -name '*.go' -a ! \( -name 'zz_generated*' -o -name '*_test.go' \))
GO_TESTS := $(shell find . -type f -name '*_test.go')
TAG_NAME = $(shell git describe --tags --abbrev=0 --exact-match 2>/dev/null)
GIT_COMMIT = $(shell git rev-parse --short=7 HEAD)
ifdef TAG_NAME
	ENVIRONMENT = production
endif
ENVIRONMENT ?= development
PREFIX = /usr/local

LD_FLAGS = -s -w -X github.com/deepsquare-io/cfctl/version.Environment=$(ENVIRONMENT) -X github.com/carlmjohnson/versioninfo.Revision=$(GIT_COMMIT) -X github.com/carlmjohnson/versioninfo.Version=$(TAG_NAME)
BUILD_FLAGS = -trimpath -a -tags "netgo,osusergo,static_build" -installsuffix netgo -ldflags "$(LD_FLAGS) -extldflags '-static'"

cfctl: $(GO_SRCS)
	go build $(BUILD_FLAGS) -o cfctl main.go

bin/cfctl-linux-x64: $(GO_SRCS)
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build $(BUILD_FLAGS) -o bin/cfctl-linux-x64 main.go

bin/cfctl-linux-arm64: $(GO_SRCS)
	GOOS=linux GOARCH=arm64 CGO_ENABLED=0 go build $(BUILD_FLAGS) -o bin/cfctl-linux-arm64 main.go

bin/cfctl-linux-arm: $(GO_SRCS)
	GOOS=linux GOARCH=arm CGO_ENABLED=0 go build $(BUILD_FLAGS) -o bin/cfctl-linux-arm main.go

bin/cfctl-win-x64.exe: $(GO_SRCS)
	GOOS=windows GOARCH=amd64 go build $(BUILD_FLAGS) -o bin/cfctl-win-x64.exe main.go

bin/cfctl-darwin-x64: $(GO_SRCS)
	GOOS=darwin GOARCH=amd64 go build $(BUILD_FLAGS) -o bin/cfctl-darwin-x64 main.go

bin/cfctl-darwin-arm64: $(GO_SRCS)
	GOOS=darwin GOARCH=arm64 go build $(BUILD_FLAGS) -o bin/cfctl-darwin-arm64 main.go

bins := cfctl-linux-x64 cfctl-linux-arm64 cfctl-linux-arm cfctl-win-x64.exe cfctl-darwin-x64 cfctl-darwin-arm64

bin/checksums.txt: $(addprefix bin/,$(bins))
	sha256sum -b $(addprefix bin/,$(bins)) | sed 's/bin\///' > $@

bin/checksums.md: bin/checksums.txt
	@echo "### SHA256 Checksums" > $@
	@echo >> $@
	@echo "\`\`\`" >> $@
	@cat $< >> $@
	@echo "\`\`\`" >> $@

.PHONY: build-all
build-all: $(addprefix bin/,$(bins)) bin/checksums.md

.PHONY: clean
clean:
	rm -rf bin/ cfctl

smoketests := smoke-basic smoke-files smoke-upgrade smoke-reset smoke-os-override smoke-init smoke-backup-restore smoke-dynamic smoke-basic-openssh smoke-dryrun
.PHONY: $(smoketests)
$(smoketests): cfctl
	$(MAKE) -C smoke-test $@

golint := $(shell which golangci-lint)
ifeq ($(golint),)
golint := $(shell go env GOPATH)/bin/golangci-lint
endif

gotest := $(shell which gotest)
ifeq ($(gotest),)
gotest := go test
endif

$(golint):
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

.PHONY: lint
lint: $(golint)
	$(golint) run ./...

.PHONY: test
test: $(GO_SRCS) $(GO_TESTS)
	$(gotest) -v ./...

.PHONY: install
install: cfctl
	install -d $(DESTDIR)$(PREFIX)/bin/
	install -m 755 cfctl $(DESTDIR)$(PREFIX)/bin/
