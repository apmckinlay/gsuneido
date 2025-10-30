<div style="float:right"><span class="builtin">Builtin</span></div>

### BuildInfo

``` suneido
() => string
```

Returns the Go debug.ReadBuildInfo. For example:

``` suneido
go	go1.20.3
path	github.com/apmckinlay/gsuneido
mod	github.com/apmckinlay/gsuneido	(devel)	
dep	github.com/google/uuid	v1.3.0	h1:t6JiXgmwXMjEs8VusXIJk2BXHsn+wx8BZdTaoZ5fu7I=
dep	github.com/kljensen/snowball	v0.8.0	h1:WU4cExxK6sNW33AiGdbn4e8RvloHrhkAssu2mVJ11kg=
dep	golang.org/x/crypto	v0.6.0	h1:qfktjS5LUO+fFKeJXZ+ikTRijMmljikvG68fpMMruSc=
dep	golang.org/x/exp	v0.0.0-20230213192124-5e25df0256eb	h1:PaBZQdo+iSDyHT053FjUCgZQ/9uqVwPOcl7KSWhKn6w=
dep	golang.org/x/sys	v0.5.0	h1:MUK/U/4lj1t1oPg0HfuXDN/Z1wv31ZJ/YcPiGccS4DU=
dep	golang.org/x/text	v0.7.0	h1:4BRB4x83lYWy72KwLD/qYDuTu7q9PjSagHvijDw7cLo=
dep	golang.org/x/time	v0.3.0	h1:rg5rLMjNzMS1RkNLzCG38eapWhnYLFYXDXj2gOlr8j4=
build	-buildmode=exe
build	-compiler=gc
build	-trimpath=true
build	CGO_ENABLED=1
build	GOARCH=amd64
build	GOOS=windows
build	GOAMD64=v1
build	vcs=git
build	vcs.revision=aec321f2ee937cb339ac01857f7a38635ed07874
build	vcs.time=2023-04-06T16:28:50Z
build	vcs.modified=true
```

See also: [Built](<Built.md>)