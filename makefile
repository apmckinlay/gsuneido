# requires sh on path (e.g. from MinGW/msys/1.0/bin)
BUILT=$(shell date "+%b %e %Y %X")

build:
	@go build -v -ldflags "-X 'main.builtDate=${BUILT}'"

test:
	go test ./...

client: build
	@./gsuneido -c t@../tok

# need 64 bit windres e.g. from mingw64
gsuneido_windows.syso : res/suneido.rc res/suneido.manifest
	windres -F pe-x86-64 -o gsuneido_windows.syso res/suneido.rc

.PHONY : build test client

# -trimpath
# -ldflags="-H windowsgui"
