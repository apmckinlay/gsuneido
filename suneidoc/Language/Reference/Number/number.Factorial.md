#### number.Factorial

``` suneido
() => number
```

Returns the factorial of the number. The factorial is obtained by multiplying every integer between 1 and the number (the factorial of 3 would be 3 * 2 * 1). Factorial does not handle negative numbers.

For example:

``` suneido
n = 3
n.Factorial()
    => 6
n = 5
n.Factorial()
    => 120
```

**Note:** The current implementation is recursive and so is limited by call nesting.