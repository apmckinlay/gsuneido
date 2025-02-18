# requires sh on path (e.g. from msys)
BUILT=$(shell date "+%b %-d %Y %R")

GO = go
GOOS = $(shell go env GOOS)
GOARCH = $(shell go env GOARCH)
OUTPUT = gs_$(GOOS)_$(GOARCH)
BUILD = build -buildvcs=true -trimpath -o $(OUTPUT)
LDFLAGS = -s -w -X 'main.builtDate=$(BUILT)'
ifdef PATHEXT
	# Windows stuff
	BUILD = build -buildvcs=true -trimpath
	OUTPUT = gsuneido.exe gsuneido.com gsport.exe
	GUIFLAGS = $(LDFLAGS) -X main.mode=gui -H windowsgui
	PORTABLE = export CGO_ENABLED=0 ; $(GO) $(BUILD) -o gsport.exe \
		-ldflags "$(LDFLAGS)" -tags portable
	CSIDE = $(GO) run cmd/deps/deps.go
endif

build:
	@$(CSIDE)
	@$(GO) version
	@rm -f $(OUTPUT)
ifdef PATHEXT
	$(GO) $(BUILD) -v -ldflags "$(GUIFLAGS)"
	$(PORTABLE)
else
	export CGO_ENABLED=0 ; $(GO) $(BUILD) -v -ldflags "$(LDFLAGS)"
endif

gsuneido:
	@$(CSIDE)
	@rm -f gsuneido.exe
	$(GO) $(BUILD) -v -ldflags "$(GUIFLAGS)"
	
# export GOEXPERIMENT=cgocheck2 ; 

race:
ifdef PATHEXT
	@$(CSIDE)
	$(GO) $(BUILD) -v -ldflags "$(GUIFLAGS)" -race -o race/
	$(PORTABLE) -race -o race/gsport.exe
else
	$(GO) $(BUILD) -v -ldflags "$(LDFLAGS)" -race -o race/$(OUTPUT)
endif


port: # a Windows version without the Windows stuff
	$(PORTABLE)
	
all: git-status build amd arm
	@mkdir -p out
	cp gsuneido.exe gsport.exe out
	mv gs_linux_amd64 gs_linux_arm64 out

# NOTE: requires test e.g. from msys
git-status:
	@test -z "$(shell git status --porcelain)"
		
arm: # linux
	export CGO_ENABLED=0 GOARCH=arm64 GOOS=linux ; $(GO) build -buildvcs=true \
		-trimpath -o gs_linux_arm64 -v -ldflags "$(LDFLAGS)"

amd: # linux
	export CGO_ENABLED=0 GOARCH=amd64 GOOS=linux ; $(GO) build -buildvcs=true \
		-trimpath -o gs_linux_amd64 -v -ldflags "$(LDFLAGS)"

test:
	export CGO_ENABLED=0 ; $(GO) test -short -vet=off -tags portable ./...

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

help:
	@echo "make [target]"
	@echo "build"
	@echo "    build for current OS"
	@echo "gsuneido"
	@echo "    build gsuneido executable"
	@echo "port"
	@echo "    build windows gsport"
	@echo "arm"
	@echo "    build arm linux executable"
	@echo "amd"
	@echo "    build amd linux executable"
	@echo "all"
	@echo "    build current OS, arm, and amd executables"
	@echo "test"
	@echo "    run tests"
	@echo "clean"
	@echo "    remove built files"

.PHONY : build gsuneido port test generate clean zap race racetest release \
    help arm amd all git-status
