#### stack.Top

``` suneido
(i = 0) => value
```

Returns the i'th top value on the stack without removing the value from the stack.
i.e. Top() returns the top value, Top(1) returns the next value down, etc.

Throws "stack underflow" if i is greater than or equal to the stack size.