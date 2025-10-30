#### string.AfterLast

``` suneido
(delimiter) => string
```

Returns the portion of the string following the last occurrence of **delimiter**.

For example:

``` suneido
str = 'first and second lines but not third and fourth lines'
str.AfterLast('and')
    => " fourth lines"
```


See also:
[string.AfterFirst](<string.AfterFirst.md>),
[string.BeforeFirst](<string.BeforeFirst.md>),
[string.BeforeLast](<string.BeforeLast.md>),
[string.LeftTrim](<string.LeftTrim.md>),
[string.RightTrim](<string.RightTrim.md>),
string.Substr,
[string.Trim](<string.Trim.md>)
