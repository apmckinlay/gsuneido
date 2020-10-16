# requires sh on path (e.g. from MinGW)
BUILT=$(shell date "+%b %e %Y %X")

LDFLAGS = -X 'main.builtDate=${BUILT}' -s -w
ifdef PATHEXT
	LDFLAGS += -H windowsgui
endif

build:
	go build -v -ldflags "$(LDFLAGS)"

all:
	go build -v -ldflags "$(LDFLAGS)" -a

console:
	go build -v -ldflags "-X 'main.builtDate=${BUILT}'"

portable:
	go build -v -ldflags "-X 'main.builtDate=${BUILT}'" -tags portable

test:
	go test -race -short -count=1 ./...

repl: build
	cmd /c start/w ./gsuneido -repl

# need the ./ so sh won't find an old one on the path
client: build
	./gsuneido.exe -c -- t@../tok

zap:
	go build -ldflags "-s -w" ./cmd/zap

generate:
	go generate -x ./...

clean:
	go clean -cache -testcache

# need 64 bit windres e.g. from mingw64
# if this fails with: 'c:\Program' is not recognized
# copy the command line and run it manually
gsuneido_windows.syso : res/suneido.rc res/suneido.manifest
	windres -F pe-x86-64 -o gsuneido_windows.syso res/suneido.rc

.PHONY : build all console portable test repl client generate clean zap

# -trimpath (but breaks vscode goto)
