PACKAGE	= i3icons2
GOPATH	= $(CURDIR)/.gopath
GOBIN	= "/usr/local/bin/"
GOFILES	= $(wildcard *.go)
CONFIG 	= $(CURDIR)/icons.config
COPY 	= cp
CONFDIR	= "/etc/"

.PHONY: all
all:    get build

.PHONY: get build 
get:
	@GOPATH=$(GOPATH) GOBIN=$(GOBIN) go get .

build:
	@echo "Building $(GOFILES)"
	@GOPATH=$(GOPATH) GOBIN=$(GOBIN) go build -o bin/$(PACKAGE) $(GOFILES)

.PHONY: install 
install:
	@GOPATH=$(GOPATH) GOBIN=$(GOBIN) go install $(GOFILES)
	${COPY} ${CONFIG} ${CONFDIR}/i3icons2.config

