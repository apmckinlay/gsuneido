# requires sh on path (e.g. from MinGW/msys/1.0/bin)
BUILT=$(shell date "+%b %e %Y %X")

build: gsuneido.syso
	@go build -v -ldflags "-X 'main.builtDate=${BUILT}'"

gsuneido.syso: suneido.manifest suneido.ico
	rsrc -manifest suneido.manifest -ico suneido.ico -o gsuneido.syso

.PHONY : build

# -ldflags="-H windowsgui"
# go get github.com/akavel/rsrc
