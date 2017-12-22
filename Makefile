
GOPATH := $(PWD)/deps/gopath
GOROOTPARENT := $(PWD)/deps/toolchain
GOROOT := $(PWD)/deps/toolchain/go
SHELL := /bin/bash # Use bash syntax

entry:
	@echo "usage:"
	@echo "make anemomteer3 : make the anemomteer executable"
	@echo "make check : checks that the algorithm compiles"
	@echo "make publish : publishes the algorithm to the server"

.PHONY: entry

deps:
	@if [[ ! -e ${GOROOT}/bin/go ]]; then \
		mkdir -p ${GOROOTPARENT}; \
		wget https://storage.googleapis.com/golang/go1.9.2.linux-amd64.tar.gz -O ${GOROOTPARENT}/go.tar.gz; \
		cd ${GOROOTPARENT} && tar -xf go.tar.gz; \
		rm ${GOROOTPARENT}/go.tar.gz; \
	fi
	${GOROOT}/bin/go get -u github.com/immesys/chirp-l7g
	${GOROOT}/bin/go get github.com/immesys/ragent

.PHONY: deps

publish: check
	@echo "Remote publishing disabled at Ed and Ali's request"

.PHONY: publish

check: deps
	cd src && ${GOROOT}/bin/go build
.PHONY: check

anemomteer3: deps
	cd src && ${GOROOT}/bin/go build && mv src ../anemomteer3
.PHONY: anemomteer3
