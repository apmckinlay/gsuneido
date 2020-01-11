# requires sh on path (e.g. from MinGW)
BUILT=$(shell date "+%b %e %Y %X")

LDFLAGS = -X 'main.builtDate=${BUILT}'
ifdef PATHEXT
	LDFLAGS += -H windowsgui
endif

build:
	go build -v -ldflags "$(LDFLAGS)"

all:
	go build -a -v -ldflags "$(LDFLAGS)"

console:
	go build -v -ldflags "-X 'main.builtDate=${BUILT}'"
	./gsuneido -repl

test:
	go test -count=1 ./...

repl: build
	cmd /c start/w ./gsuneido -repl

# need the ./ so sh won't find an old one on the path
client: build
	./gsuneido.exe -c -- t@../tok

generate:
	go generate ./...

# need 64 bit windres e.g. from mingw64
gsuneido_windows.syso : res/suneido.rc res/suneido.manifest
	windres -F pe-x86-64 -o gsuneido_windows.syso res/suneido.rc

.PHONY : build all console test repl client generate

# -trimpath (but breaks vscode goto)
