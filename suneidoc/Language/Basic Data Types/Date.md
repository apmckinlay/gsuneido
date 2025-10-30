### Date

Suneido has a built-in Date type that contains both a date and a time, with millisecond precision.

Date literals are written as one of:

``` suneido
#yyyyMMdd
#yyyyMMdd.HHmm
#yyyyMMdd.HHmmss
#yyyyMMdd.HHmmssttt // where ttt is milliseconds
```

Time elements that are not supplied are set to zero.

Dates can be constructed more flexibly using the [Date constructor](<../Reference/Date/Date.md>).

See also:
[Date](<../Reference/Date.md>)