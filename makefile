# requires sh on path (e.g. from MinGW)
BUILT=$(shell date "+%b %e %Y %X")

EXE = gsuneido
LDFLAGS = -s -w -X 'main.builtDate=${BUILT}'
GUIFLAGS = $(LDFLAGS)
ifdef PATHEXT
	# Windows stuff
	EXE = gsuneido.exe gsuneido.com gsport.exe
	GUIFLAGS = $(LDFLAGS) -X main.mode=gui -H windowsgui
	CONSOLE = go build -o gsuneido.com -ldflags "$(LDFLAGS)" -tags com
	PORTABLE = export CGO_ENABLED=0 ; go build -o gsport.exe -ldflags "$(LDFLAGS)" -tags portable
endif

build:
	@rm -f $(EXE)
	@go version
	go build -v -ldflags "$(GUIFLAGS)"
	$(CONSOLE)
	$(PORTABLE)

race:
	go build -v -ldflags "$(GUIFLAGS)" -race

portable:
	# a Windows version without the Windows stuff
	$(PORTABLE)

test:
	go test -short ./...

racetest:
	go test -race -short -count=1 ./...

zap:
	go build -ldflags "-s -w" ./cmd/zap

generate:
	go generate -x ./...

clean:
	rm -f $(EXE)
	go clean -cache -testcache

# need 64 bit windres e.g. from mingw64
# if this fails with: 'c:\Program' is not recognized
# copy the command line and run it manually
gsuneido_windows.syso : res/suneido.rc res/suneido.manifest
	windres -F pe-x86-64 -o gsuneido_windows.syso res/suneido.rc

.PHONY : build portable test generate clean zap race

# -trimpath (but breaks vscode goto)
