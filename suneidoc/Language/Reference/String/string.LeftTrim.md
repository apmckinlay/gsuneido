#### string.LeftTrim

``` suneido
(chars = " \t\r\n") => string
```

Returns a copy of the string with leading characters removed. The characters to be removed default to space, tab, return, and linefeed.

For example:

``` suneido
"   hello world   ".LeftTrim()
    => "hello world   "
```


See also:
[string.AfterFirst](<string.AfterFirst.md>),
[string.AfterLast](<string.AfterLast.md>),
[string.BeforeFirst](<string.BeforeFirst.md>),
[string.BeforeLast](<string.BeforeLast.md>),
[string.RightTrim](<string.RightTrim.md>),
string.Substr,
[string.Trim](<string.Trim.md>)
