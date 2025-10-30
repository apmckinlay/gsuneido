### Implicit Dynamic Parameters

Parameters that start with an underscore can receive an implicit value from [Dynamic Variables](<../Names/Dynamic Variables.md>).

For example:

``` suneido
f = function (_context) { doSomthingWith(context) }
```

Within the function, the name is used without the underscore. If you wrote _name within the function body you would be referring to the dynamic variable, not the parameter.

context can be passed as a normal argument:

``` suneido
f(myContext)
f(context: myContext)
```

Or it can receive context via a dynamic variable:

``` suneido
_context = myContext
...
f()
```

The dynamic variable does not need to be set in the same function as the call, it can be set in anywhere in the sequence of calls leading here.

Default values can still be supplied. The default value is only used if no argument is passed and no dynamic variable is found.

``` suneido
f = function (_p = 0) { return p }
f()
    => 0
f(123)
    => 123
_p = 456
f()
    => 456
f(123)
    => 123
```