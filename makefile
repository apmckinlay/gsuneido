# requires sh on path (e.g. from MinGW/msys/1.0/bin)
BUILT=$(shell date "+%b %e %Y %X")

build:
	@go build -race -ldflags "-X 'main.builtDate=${BUILT}'" gsuneido.go

.PHONY : build
