### Nothing

``` suneido
()
```

An empty function with no return value.

Primarily useful when writing a block that sometimes returns nothing. For example:

``` suneido
{ it < 0 ? Nothing() : it.Sqrt() }
```

It can also used to make it more obvious that a function is not returning anything, for example:

``` suneido
return Nothing()
```