### ReadableDuration

``` suneido
(seconds) => string
.Between(t1, t2) => string
```

Returns a human readable string representation of the difference between two times. Only 3 digits of precision are kept. Possible units are hr, min, sec, and ms.

For example:

``` suneido
ReadableDuration(123)
=> "123 ms"

ReadableDuration(12.34)
=> "12.3 sec"

ReadableDuration(7200)
=> "2 hrs"
```

See also:
[Date](<Date.md>),
[ReadableSize](<ReadableSize.md>),
[Stopwatch](<Stopwatch.md>)