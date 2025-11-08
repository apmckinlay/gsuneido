# The native non-gui version is pure Go.
# It should build with just Go on Windows, Linux, and Mac.
# The windows amd64 gui version is built with mingw-w64 on the path
# The windows arm64 gui version is built with llvm-mingw on Mac.
# It is NOT fully functional.
# The windows amd64 gui version will run on Arm Windows with emulation.
# requires sh (and date and rm) on path (e.g. from msys2)

BUILT=$(shell date "+%b %-d %Y %R")

GOOS = $(shell go env GOOS)
GOARCH = $(shell go env GOARCH)
LDFLAGS = -s -w -X 'main.builtDate=$(BUILT)'
BUILD = build -v -buildvcs=true -trimpath

EXE =
ifdef PATHEXT
EXE = .exe
endif

# build compiles the native non-gui version
build:
	@go version
	CGO_ENABLED=0 \
	go $(BUILD) -o gs_$(GOOS)_$(GOARCH)$(EXE) \
	  -ldflags "$(LDFLAGS)" 

define BUILD_BINARY
	$(eval GO_VARS := $(subst _, ,$*))
	$(eval GOOS := $(word 1,$(GO_VARS)))
	$(eval GOARCH := $(word 2,$(GO_VARS)))
	CGO_ENABLED=0 \
	GOARCH=$(GOARCH) GOOS=$(GOOS) \
	go $(BUILD) -o $@ -ldflags "$(LDFLAGS)"
endef

gs_%: FORCE
	$(call BUILD_BINARY)

gs_%.exe: FORCE
	$(call BUILD_BINARY)

WINGUI = -X main.mode=gui -H windowsgui
	
gs_windows_amd64_gui.exe: FORCE gsuneido_windows_amd64.syso
	@go version
	go run cmd/deps/deps.go
	CGO_ENABLED=1 \
	go $(BUILD) -tags gui -ldflags "$(LDFLAGS) $(WINGUI)" -o $@

gsuneido_windows_amd64.syso : res/suneido.rc res/suneido.manifest
	windres -F pe-x86-64 -o gsuneido_windows_amd64.syso res/suneido.rc

gui: gs_windows_amd64_gui.exe

both: build gui
	
deploy: git-status windows_amd64.exe windows_amd64_gui gs_linux_arm64 gs_linux_amd64
	@mkdir -p deploy
	cp gs_windows_amd64.exe deploy\gsport.exe
	cp gs_windows_arm64_gui.exe deploy\gsuneido.exe
	mv gs_linux_amd64 gs_linux_arm64 deploy

# NOTE: requires test e.g. from msys
git-status:
	@test -z "$(shell git status --porcelain)"

test:
	CGO_ENABLED=0 \
	go test -short -vet=off -timeout 30s ./...

racetest:
	go test -race -short -count=1 ./...

zap:
	go build -ldflags "-s -w" ./cmd/zap

generate:
	go generate -x ./...

clean:
	go clean -cache -testcache

# for cross compiling on Arm Mac for Arm Windows
LLVM_MINGW = /Users/andrew/apps/llvm-mingw/bin/aarch64-w64-mingw32

gs_windows_arm64_gui.exe: FORCE gsuneido_windows_arm64.syso
	CGO_ENABLED=1 \
	GOARCH=arm64 GOOS=windows \
	CC=$(LLVM_MINGW)-clang \
	CXX=$(LLVM_MINGW)-clang++ \
	CGO_CXXFLAGS=-Wno-inconsistent-missing-override \
	go $(BUILD) -tags gui -o gs_windows_arm64_gui.exe \
	  -ldflags "$(LDFLAGS) $(WINGUI)"

gsuneido_windows_arm64.syso : res/suneido.rc res/suneido.manifest
	$(LLVM_MINGW)-windres -o gsuneido_windows_arm64.syso res/suneido.rc

release:
	./gsuneido -dump stdlib
	./gsuneido -dump suneidoc
	./gsuneido -dump imagebook
	-mkdir release
	mv stdlib.su suneidoc.su imagebook.su release
	cd release && \
	  rm suneido.db ; \
	  ../gsuneido -load stdlib && \
	  ../gsuneido -load suneidoc && \
	  ../gsuneido -load imagebook

help:
	@echo "make [target]"
	@echo "build"
	@echo "    compile a native non-gui version"
	@echo "gs_<GOOS>_<GOARCH>[.exe]"
	@echo "    compile windows/linux/darwin amd64/arm64 version"
	@echo "deploy"
	@echo "    build and copy to deploy directory"
	@echo "    windows_amd64.exe windows_amd64_gui gs_linux_arm64 gs_linux_amd64"
	@echo "test"
	@echo "    run tests"
	@echo "clean"
	@echo "    remove built files"

.PHONY : build test generate clean zap race racetest release \
    help deploy git-status both gui FORCE
