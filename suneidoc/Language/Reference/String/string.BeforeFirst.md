#### string.BeforeFirst

``` suneido
(delimiter) => string
```

Returns a string containing the portion that comes before the first occurrence of **delimiter**.

For example:

``` suneido
str = 'first and second lines but not third and forth lines'
str.BeforeFirst('lines')
    => "first and second "
```


See also:
[string.AfterFirst](<string.AfterFirst.md>),
[string.AfterLast](<string.AfterLast.md>),
[string.BeforeLast](<string.BeforeLast.md>),
[string.LeftTrim](<string.LeftTrim.md>),
[string.RightTrim](<string.RightTrim.md>),
string.Substr,
[string.Trim](<string.Trim.md>)
