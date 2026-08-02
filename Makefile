GO ?= go

build:
	$(GO) build -o bin/gvte ./cmd

run:
	$(GO) run bin/gvte
