<div style="float:right"><span class="builtin">Builtin</span></div>

#### file.Flush

Ensures that buffered output has been written to the file.

For normal usage, it is not necessary to explicitly call Flush.

Note: Close automatically flushes.

For example:

``` suneido
File("tmp", "w")
    {|f|
    f.Writeline("some stuff")
    f.Flush()
    ...
    f.Writeline("more stuff")
    }
```