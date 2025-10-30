#### string.RemovePrefix

``` suneido
(prefix) => string
```

Returns the string without the specified prefix. If the string does not have the specified prefix it will be returned unchanged.

For example:

``` suneido
"foobar".RemovePrefix("foo")
=> "bar"
```

See also: [string.Prefix?](<string.Prefix?.md>), [string.RemoveSuffix](<string.RemoveSuffix.md>)