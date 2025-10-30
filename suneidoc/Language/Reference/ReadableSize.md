### ReadableSize

``` suneido
(n) => string
ReadableSize.FromInt(n) => string
ReadableSize.ToInt(s) => n
```

Converts to and from a human readable string representation of a size. Converting to a string only keeps 4 digits of precision.

For example:

``` suneido
ReadableSize(100000)
    => "97.6kb"

ReadableSize.ToInt("97.6kb")
    => 99942
```

See also:
[ReadableDuration](<ReadableDuration.md>)