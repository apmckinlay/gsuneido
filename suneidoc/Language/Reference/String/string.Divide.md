#### string.Divide

``` suneido
(n = 1) => object
```

Returns an object containing the string split into substrings of length n.

For example:

``` suneido
"hello".Divide(2)
    => #("he", "ll", "o")
```

The last substring may be less than n characters.

If the string is "" then an empty object is returned.

See also: [string.MapN](<string.MapN.md>)