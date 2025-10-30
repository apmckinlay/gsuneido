### Chaining Method Calls

``` suneido
object.
    method1(...).
    method2(...).
    method3(...)
```

Note: The period must be at the end of the line so that Suneido knows the line is continued.

For example:

``` suneido
    where = ShortestKey(query).
        Split(',').
        Map!({|f| f $ ' = ' $ Display(record[fld]) }).
        Join(' and ')
```