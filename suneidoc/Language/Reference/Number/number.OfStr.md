#### number.OfStr

``` suneido
(block) => string
```

Returns the concatenation of multiple calls to block.

For example:

``` suneido
4.OfStr({ '12345'.RandChar() })
=> "2245"
```

For multiple copies of a string, use [string.Repeat](<../String/string.Repeat.md>)

See also: [number.Of](<number.Of.md>)