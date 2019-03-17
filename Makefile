#
# github.com/scott-haines/sinject
#
all: binary

export VERSION = $(shell cat sinject.VERSION)

binary: ## build executable for Linux
	./scripts/build/binary