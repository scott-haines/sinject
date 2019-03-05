#
# github.com/scott-haines/sinject
#
all: binary

export VERSION = $(shell cat VERSION)

binary: ## build executable for Linux
	./scripts/build/binary