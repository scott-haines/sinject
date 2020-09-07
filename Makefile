#
# github.com/scott-haines/sinject
#
all: clean binaries

export VERSION = $(shell cat sinject.VERSION)

.PHONY: clean
clean: ## clean the build directory
	./scripts/build/clean

ARCHS = darwin linux
.PHONY: binaries
binaries: ## build executables
	$(foreach arch,$(ARCHS),env GOOS=$(arch) ./scripts/build/binary;)

install: ## compile and install locally
	env GOOS=linux ./scripts/build/binary
	mkdir -p $$GOPATH/bin/
	cp build/sinject-linux-amd64 $$GOPATH/bin/sinject