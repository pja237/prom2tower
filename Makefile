.PHONY: clean build test test_new list all

.ONESHELL:
#SHELL =
#.SHELLFLAGS =

# Inject into binary via linker:
# ...in github actions comes from make -e version=git_ref
name=prom2tower
version=$(shell cat VERSION)
commit=$(shell git show --no-patch --format=format:%H HEAD)
buildVersionVar=github.com/pja237/prom2tower/internal/version.buildVersion
buildCommitVar=github.com/pja237/prom2tower/internal/version.buildCommit

# various directories
bindirs=$(wildcard ./cmd/*)
installdir=build/$(name)-$(version)

# list of files to include in build
bins=$(notdir $(bindirs))
readme=README.md
templates=
config=cmd/prom2tower/conf.yaml.example

# can be replaced with go test ./... construct
testdirs=$(sort $(dir $(shell find ./ -name *_test.go)))

all: list test build install

list:
	@echo "================================================================================"
	@echo "bindirs  found: $(bindirs)"
	@echo "bins     found: $(bins)"
	@echo "testdirs found: $(testdirs)"
	@echo "================================================================================"

build:
	@echo "********************************************************************************"
	@echo Building $(bindirs)
	@echo Variables:
	@echo buildVersionVar: $(buildVersionVar)
	@echo version: $(version)
	@echo buildCommitVar: $(buildCommitVar)
	@echo commit: $(commit)
	@echo "********************************************************************************"
	for i in $(bindirs);
	do
		echo "................................................................................"
		echo "--> Now building: $$i"
		echo "................................................................................"
		go build -v -ldflags '-X $(buildVersionVar)=$(version) -X $(buildCommitVar)=$(commit)' $$i;
	done;

install:
	mkdir -p $(installdir)
	cp $(bins) $(readme) $(templates) $(config) $(installdir)

test_new:
	$(foreach dir, $(testdirs), go test -v -count=1 $(dir) || exit $$?;)

test:
	@echo "********************************************************************************"
	@echo Testing
	@echo "********************************************************************************"
	go test -v -count=1 ./...

clean:
	rm $(bins)
	rm -rf $(installdir)
