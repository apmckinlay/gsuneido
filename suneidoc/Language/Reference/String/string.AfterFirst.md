#### string.AfterFirst

``` suneido
(delimiter) => string
```

Returns a string containing the portion that comes after the first occurrence of **delimiter**.

For example:

``` suneido
str = 'first and second lines but not third and fourth lines'
str.AfterFirst('lines')
    => " but not third and fourth lines"
```


See also:
[string.AfterLast](<string.AfterLast.md>),
[string.BeforeFirst](<string.BeforeFirst.md>),
[string.BeforeLast](<string.BeforeLast.md>),
[string.LeftTrim](<string.LeftTrim.md>),
[string.RightTrim](<string.RightTrim.md>),
string.Substr,
[string.Trim](<string.Trim.md>)
