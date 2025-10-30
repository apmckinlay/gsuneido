#### string.Trim

``` suneido
(chars = " \t\r\n") => string
```

Returns a copy of the string with leading and trailing characters removed. The characters to be removed default to space, tab, return, and linefeed.

For example:

``` suneido
"   hello world   ".Trim()
    => "hello world"
```


See also:
[string.AfterFirst](<string.AfterFirst.md>),
[string.AfterLast](<string.AfterLast.md>),
[string.BeforeFirst](<string.BeforeFirst.md>),
[string.BeforeLast](<string.BeforeLast.md>),
[string.LeftTrim](<string.LeftTrim.md>),
[string.RightTrim](<string.RightTrim.md>),
string.Substr
