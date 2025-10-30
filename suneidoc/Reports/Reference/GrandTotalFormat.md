### GrandTotalFormat

``` suneido
( format )
```

Displays the format with a double horizontal bar over it, i.e. like a grand total.  For example:

``` suneido
( Query
    ( _output desc: "Total", amount: ( GrandTotal total_amount ) )
```

See also:
[TotalFormat](<TotalFormat.md>)