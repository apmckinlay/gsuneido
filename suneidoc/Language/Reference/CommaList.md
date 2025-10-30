### CommaList

``` suneido
(@args) => string
```

Combines the arguments into a comma separated string, ignoring any arguments that are empty strings.

For example:

``` suneido
first = "Andrew"
middle = ""
last = "McKinlay"
CommaList(last, first, middle)
    => "McKinlay, Andrew"
```