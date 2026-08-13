Building gSuneido
================

The easiest way to build is to use `make`.
`go build` is not sufficient, especially on Windows.

The non-gui version is pure Go and should only require Go to build.

The Windows GUI version needs [cgo](https://golang.org/cmd/cgo/) so it requires a C/C++ compiler on the path. [w64devkit](https://github.com/skeeto/w64devkit) is a good choice and includes make and sh.

I use the included makefile (requires make). If you prefer not to use make you can just look at the makefile to see what the build commands are.

You should be able to run the tests with the usual: `go test -short ./..`
Or you can use `make test`

I try to fix warnings but depending on what checker you use, you may get some.

My normal development environment is [Visual Studio Code](https://code.visualstudio.com/) with the [Go extension](https://marketplace.visualstudio.com/items?itemName=ms-vscode.Go).

Building liblexilla.a and libscintilla.a
----------------------------------------

w64devkit works to build these, following their instructions

To build a smaller version of liblexilla with specific lexers:
(as of lexilla 5.4.5)
- remove the unneeded lexers from the lexers directory
- the ones used by stdlib are:
  - LexCPP.cxx (for JavaScript)
  - LexCSS.cxx
  - LexHTML.cxx
  - LexMarkdown.cxx
- from the scripts directory, run: python LexillaGen.py 
- rebuild (if you built with the full set, you need to run make clean first)
