test-all: test-124 test-125
GOFLAGS=-ldflags=-checklinkname=0
export GOFLAGS

install-go:
	@which go1.24.6 || (go install golang.org/dl/go1.24.6@latest && go1.24.6 download)
	@which go1.25.1 || (go install golang.org/dl/go1.25.1@latest && go1.25.1 download)

test-124:
	go1.24.6 test -count=1 ./...

test-125:
	go1.25.1 test -count=1 ./...
