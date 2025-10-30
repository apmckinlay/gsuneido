### Ensuring Execution on the Server

Sometimes you want to ensure that a function or method is executed on the server (using [ServerEval](<../../Language/Reference/ServerEval.md>)).

Requiring the callers to use ServerEval is error prone and sooner or later someone will forget.

One way to handle this is to have an extra method:

``` suneido
MyClass

class
    {
    Func(a, b, c)
        {
        ServerEval('MyClass.Func2', a, b, c)
        }
    Func2(a, b, c)
        {
        ...
        }
```

However, this requires two methods, and Func2 has to be public which means someone could call it directly by mistake.

A simpler, cleaner way to do it is to use the same function and detect whether you are running as a client (with [Client?](<../../Language/Reference/Client?.md>)()):

``` suneido
MyFunc

function (a, b, c)
    {
    if Client?()
        ServerEval('MyFunc', a, b, c)
    else
        {
        ...
        }
    }
```

(This example is a standalone function, but the same approach can be used with a class method.)