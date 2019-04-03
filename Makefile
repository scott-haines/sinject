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