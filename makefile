BUILT=$(shell date "+%b %e %Y %X")

build:
	go build -ldflags "-X 'main.builtDate=${BUILT}'" gsuneido.go

.PHONY : build
