<div style="float:right"><span class="builtin">Builtin</span></div>

### LONG

Can be used to convert Suneido numbers to and from "native" long integer representation.

For example:

``` suneido
LONG(Object(x: 64))

    => "@\x00\x00\x00"

LONG("@\x00\x00\x00")

    => #(x: 64)
```