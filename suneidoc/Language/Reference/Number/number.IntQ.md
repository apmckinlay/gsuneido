#### number.Int?

``` suneido
() => number
```

Returns true if the number is an integer, and false if it has a fractional component.

For example:

``` suneido
n = 123
n.Int?()
    => true
n += 0.456
n.Int?()
    => false
n.Int().Int?()
    => true
```