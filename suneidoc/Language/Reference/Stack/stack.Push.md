#### stack.Push

``` suneido
(value)
(@members)
```

Pushes a new value onto the top of the stack.

Can be used to push a single value:

``` suneido
stack.Push(123);
```

or an object created from the arguments:

``` suneido
stack.Push(name: "Fred", age: 25)
stack.Pop() => #(name: "Fred", age: 25)
```