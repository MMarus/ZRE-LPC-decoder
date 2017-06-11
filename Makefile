GOPATH=$(shell pwd)

all: build 
build:
	@echo "Building in $(GOPATH)"
	@GOPATH=$(GOPATH) go build src/zre_proj1_linux.go
