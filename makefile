# requires sh on path (e.g. from MinGW)
BUILT=$(shell date "+%b %e %Y %X")

GO = go1.18rc1
EXE = gsuneido
LDFLAGS = -s -w -X 'main.builtDate=${BUILT}'
GUIFLAGS = $(LDFLAGS)
ifdef PATHEXT
	# Windows stuff
	EXE = gsuneido.exe gsuneido.com gsport.exe
	GUIFLAGS = $(LDFLAGS) -X main.mode=gui -H windowsgui
	CONSOLE = $(GO) build -o gsuneido.com -ldflags "$(LDFLAGS)" -tags com
	PORTABLE = export CGO_ENABLED=0 ; $(GO) build -o gsport.exe -ldflags "$(LDFLAGS)" -tags portable
endif

build:
	@rm -f $(EXE)
	@$(GO) version
	$(GO) build -v -ldflags "$(GUIFLAGS)"
	$(CONSOLE)
	$(PORTABLE)

race:
	$(GO) build -v -ldflags "$(GUIFLAGS)" -race

portable:
	# a Windows version without the Windows stuff
	$(PORTABLE)

test:
	$(GO) test -short ./...

racetest:
	$(GO) test -race -short -count=1 ./...

zap:
	$(GO) build -ldflags "-s -w" ./cmd/zap

generate:
	$(GO) generate -x ./...

clean:
	rm -f $(EXE)
	$(GO) clean -cache -testcache

# need 64 bit windres e.g. from mingw64
# if this fails with: 'c:\Program' is not recognized
# copy the command line and run it manually
gsuneido_windows.syso : res/suneido.rc res/suneido.manifest
	windres -F pe-x86-64 -o gsuneido_windows.syso res/suneido.rc

.PHONY : build portable test generate clean zap race

# -trimpath (but breaks vscode goto)
