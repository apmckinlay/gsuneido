# requires sh on path (e.g. from MinGW)
BUILT=$(shell date "+%b %-d %Y %R")

GO = go
BUILD = build -buildvcs=true
EXE = gsuneido
LDFLAGS = -s -w -X 'main.builtDate=${BUILT}'
GUIFLAGS = $(LDFLAGS)
RACEEXE = gsrace
ifdef PATHEXT
	# Windows stuff
	BUILD = build -buildvcs=true -trimpath
	EXE = gsuneido.exe gsuneido.com gsport.exe
	GUIFLAGS = $(LDFLAGS) -X main.mode=gui -H windowsgui
	CONSOLE = $(GO) $(BUILD) -o gsuneido.com -ldflags "$(LDFLAGS)" -tags com
	PORTABLE = export CGO_ENABLED=0 ; $(GO) $(BUILD) -o gsport.exe -ldflags "$(LDFLAGS)" -tags portable
	RACEEXE = gsrace.exe
endif

build:
	@rm -f $(EXE)
	@$(GO) version
	$(GO) $(BUILD) -v -ldflags "$(GUIFLAGS)"
	$(CONSOLE)
	$(PORTABLE)

gsuneido:
	$(GO) $(BUILD) -v -ldflags "$(GUIFLAGS)"

race:
	$(GO) $(BUILD) -v -o $(RACEEXE) -ldflags "$(GUIFLAGS)" -race

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

.PHONY : build gsuneido portable test generate clean zap race racetest

# -trimpath (but breaks vscode goto)
