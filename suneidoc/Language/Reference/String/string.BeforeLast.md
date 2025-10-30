#### string.BeforeLast

``` suneido
(delimiter) => string
```

Returns the portion of the string preceding the last occurrence of **delimiter**.

For example:

``` suneido
str = 'first and second lines but not third and forth lines'
str.BeforeLast('and')
    => "first and second lines but not third "
```


See also:
[string.AfterFirst](<string.AfterFirst.md>),
[string.AfterLast](<string.AfterLast.md>),
[string.BeforeFirst](<string.BeforeFirst.md>),
[string.LeftTrim](<string.LeftTrim.md>),
[string.RightTrim](<string.RightTrim.md>),
string.Substr,
[string.Trim](<string.Trim.md>)
