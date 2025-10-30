### Pipe Operator

pipe-expression:

``` suneido
expression |> callable
```

The pipe operator `|>` provides a convenient way to chain function calls by passing the left operand as the argument to the right operand. It transforms:

``` suneido
value |> callable
```

into:

``` suneido
callable(value)
```

For example:

``` suneido
value |> Func1 |> Func2 |> Func3
```

is equivalent to:

``` suneido
Func3(Func2(Func1(value)))
```

The benefit is that the functions are written in the order they will be executed (1,2,3).

#### Precedence

The pipe operator has the lowest precedence. For example:

``` suneido
a ? b : c |> Func
```

is equivalent to:

``` suneido
Func(a ? b : c)
```

#### Performance

The pipe operator is converted to regular function calls at compile time. 
So the execution speed is exactly the same.

See also [Function Calls](<Function Calls.html>)