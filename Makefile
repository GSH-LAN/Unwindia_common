PROJECT := $(shell git config --local remote.origin.url|sed -n 's#.*/\([^.]*\)\.git#\1#p'|sed 's/[A-Z]/\L&/g')

.PHONY: \
	go-build \
	go-test

test: \
	go-test

go-build:
	go build ./...

go-test:
	go test -v -count=1 -race ./...
