#### object.UniqueValues

``` suneido
() => object
```

Returns a list of the unique values from the object.

For example:

``` suneido
#(6, 2, 7, 2, 4, 6, 5).UniqueValues()
    => #(6, 2, 7, 4, 5)
```

Preserves the order of the values.

**Note:** Member names are ignored and will not be carried over to the result.

**Warning**: May be slow for large number of unique values.

Faster alternatives may be:

``` suneido
object.Sort!().Unique!()
```

Or if you do the unique'ing as you collect the values, and take advantage of the fast hash access to members:

``` suneido
ob = Object()
ob[value] = true
...
ob = ob.Members()
```

**Note**: These alternate approaches do not preserve the order of the values.

See also:
[string.UniqueChars](<../String/string.UniqueChars.md>)