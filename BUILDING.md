Building gSuneido
================

The easiest way to build is to use `make`.
`go build` is not sufficient, especially on Windows.

The non-gui version is pure Go and should only require Go to build.

The Windows GUI version needs [cgo](https://golang.org/cmd/cgo/) so it requires a C/C++ compiler on the path. I am using [Mingw-w64](https://sourceforge.net/projects/mingw-w64/). You can either install Msys2 or just Mingw-w64. The specific version of GCC shouldn't be critical (but not too old). I would recommend using a normal install of Go, **not** installing it in Msys2 with pacman. Go just needs gcc on the path.

From the terminal, running:

    g++ --version

Should give something like:

    g++ (x86_64-posix-seh-rev0, Built by MinGW-W64 project) 8.1.0

I use the included makefile (requires make). If you prefer not to use make you can just look at the makefile to see what the build commands are. There is a make included with mingw_w64 but it's called `mingw32_make.exe` You can make a copy and rename it to `make.exe`

You should be able to run the tests with the usual: `go test -short ./..`
Or you can use `make test`

To run the [portable tests](https://github.com/apmckinlay/suneido_tests) (shared with cSuneido and jSuneido) you need to clone or download the separate Github repository. It needs to be a peer of the gsuneido directory, and must be called suneido_tests. e.g. `c:/stuff/mygsuneido` and `c:/stuff/suneido_tests`

I try to fix warnings but depending on what checker you use, you may get some.

My normal development environment is [Visual Studio Code](https://code.visualstudio.com/) with the [Go extension](https://marketplace.visualstudio.com/items?itemName=ms-vscode.Go).

Building liblexilla.a and libscintilla.a
----------------------------------------

mingw-w64 works to build these, following their instructions

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
