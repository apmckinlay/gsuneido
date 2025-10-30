<div style="float:right"><span class="builtin">Builtin</span></div>

#### file.Size

``` suneido
() => size
```

New As of BuiltDate 20241216

Returns the current size of the file.

Note: If open for "w" or "a" this flushes the file.
For best performance avoiding doing this too often while writing.

See also: [FileSize](<../FileSize.md>)