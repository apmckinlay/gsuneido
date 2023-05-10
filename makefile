# requires sh on path (e.g. from MinGW)
BUILT=$(shell date "+%b %-d %Y %R")

GO = go
GOOS = $(shell go env GOOS)
GOARCH = $(shell go env GOARCH)
ifneq ($(GOOS),darwin)
TRIMPATH = -trimpath
endif
OUTPUT = gs_$(GOOS)_$(GOARCH)
BUILD = build -buildvcs=true $(TRIMPATH) -o $(OUTPUT)
LDFLAGS = -s -w -X 'main.builtDate=$(BUILT)'
ifdef PATHEXT
	# Windows stuff
	BUILD = build -buildvcs=true -trimpath
	OUTPUT = gsuneido.exe gsuneido.com gsport.exe
	GUIFLAGS = $(LDFLAGS) -X main.mode=gui -H windowsgui
	CONSOLE = $(GO) $(BUILD) -o gsuneido.com -ldflags "$(LDFLAGS)" -tags com
	PORTABLE = export CGO_ENABLED=0 ; $(GO) $(BUILD) -o gsport.exe -ldflags "$(LDFLAGS)" -tags portable
endif

help:
	@echo "make [target]"
	@echo "build"
	@echo "    build gsuneido"
	@echo "gsuneido"
	@echo "    build gsuneido executable"
	@echo "test"
	@echo "    run tests"
	@echo "clean"
	@echo "    remove built files"

build:
	@$(GO) version
	@rm -f $(OUTPUT)
ifdef PATHEXT
	$(GO) $(BUILD) -v -ldflags "$(GUIFLAGS)"
	$(CONSOLE)
	$(PORTABLE)
else
	export CGO_ENABLED=0 ; $(GO) $(BUILD) -v -ldflags "$(LDFLAGS)"
endif

gsuneido:
	@rm -f gsuneido.exe
	$(GO) $(BUILD) -v -ldflags "$(GUIFLAGS)"

race:
ifdef PATHEXT
	$(GO) $(BUILD) -v -ldflags "$(GUIFLAGS)" -race -o race/
	$(PORTABLE) -race -o race/gsport.exe
else
	$(GO) $(BUILD) -v -ldflags "$(GUIFLAGS)" -race -o race/$(OUTPUT)
endif

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
	rm -f $(OUTPUT)
	$(GO) clean -cache -testcache

# need 64 bit windres e.g. from mingw64
# if this fails with: 'c:\Program' is not recognized
# copy the command line and run it manually
gsuneido_windows.syso : res/suneido.rc res/suneido.manifest
	windres -F pe-x86-64 -o gsuneido_windows.syso res/suneido.rc

release:
	./gsuneido -dump stdlib
	./gsuneido -dump suneidoc
	./gsuneido -dump imagebook
	-mkdir release
	cp stdlib.su suneidoc.su imagebook.su release
	cd release && \
	  rm suneido.db ; \
	  ../gsuneido -load stdlib && \
	  ../gsuneido -load suneidoc && \
	  ../gsuneido -load imagebook

.PHONY : build gsuneido portable test generate clean zap race racetest release
