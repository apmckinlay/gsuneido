<div style="float:right"><span class="builtin">Builtin</span></div>

### Display

``` suneido
(value, quotes = 0) => string
```

Returns a string version of any value.
Where reasonable, the resulting string can be evaluated to get an equal value back. However, Display is intended for human use; it is NOT intended as a way of converting data to and from strings.

The quotes argument can be:

0 = default behavior - single-quotes for single characters, double-quotes unless single or back-quotes would require less escaping   

1 = single-quotes   

2 = double-quotes

All non-printable characters will be escaped, including tab (\t), return (\r), and linefeed (\n).