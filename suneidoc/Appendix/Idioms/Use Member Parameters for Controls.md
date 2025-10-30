### Use Member Parameters for Controls

When a user interface is fixed, you can write something like:

``` suneido
Controller
    {
    Controls: (...)
```

Or if you need to build the layout in code:

``` suneido
Controller
    {
    Controls()
        {
        return ...
        }
```

When a user interface layout depends on arguments, one way to handle it is:

``` suneido
Controller
    {
    New(a, b, c)
        {
        super(.controls(a, b, c))
        }
    controls(a, b, c)
        {
        return Object( ... )
        }
```

This works but it is somewhat verbose and you have to repeat the arguments multiple times. This can be simplified with member parameters.

``` suneido
Controller
    {
    New(.a, .b, .c)
        {
        }
    Controls()
        {
        // can use .a, .b, and .c here
        return ...
        }
```

You still need to have a `New` since it must accept the arguments, but it can be empty. (An implicit [super](<../../Language/Classes/super.md>) call will be generated.)

Note: You could not use this approach prior to member parameters because an explicit super call must be the first line of New.

<div style="color: red;">

``` suneido
    New(a, b, c)
        {
        .a = a
        .b = b
        .c = c
        super() // INVALID BECAUSE NOT FIRST STATEMENT
        }
```

</div>