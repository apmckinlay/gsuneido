# requires sh on path (e.g. from MinGW/msys/1.0/bin)
BUILT=$(shell date "+%b %e %Y %X")

build: gsuneido.syso
	@go build -v -ldflags "-X 'main.builtDate=${BUILT}'"

# need 64 bit windres e.g. from mingw64
gsuneido.syso : res/suneido.rc res/suneido.manifest
	windres -F pe-x86-64 -o gsuneido.syso res/suneido.rc

.PHONY : build

# -ldflags="-H windowsgui"
