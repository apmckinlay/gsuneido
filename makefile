# requires sh on path (e.g. from MinGW)
BUILT=$(shell date "+%b %e %Y %X")

build:
	go build -v -ldflags "-X 'main.builtDate=${BUILT}' -H windowsgui"

console:
	go build -v -ldflags "-X 'main.builtDate=${BUILT}'"

test:
	go test -count=1 ./...

repl: build
	cmd /c start/w ./gsuneido -repl

client: build
	gsuneido -c t@../tok

generate:
	go generate ./...

# need 64 bit windres e.g. from mingw64
gsuneido_windows.syso : res/suneido.rc res/suneido.manifest
	windres -F pe-x86-64 -o gsuneido_windows.syso res/suneido.rc

.PHONY : build test client

# -trimpath (but breaks vscode goto)
