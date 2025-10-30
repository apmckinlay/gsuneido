### Use Client? to Run on Server

Instead of having a second function just to ServerEval you can do:

``` suneido
MyFunc

function (arg)
    {
    if Client?()
        return ServerEval("MyFunc", arg)
    ...
    return ...
    }
```

You can also use this for a class (static) method:

``` suneido
MyClass

class
    {
    MyMethod(arg)
        {
        if Client?()
            return ServerEval("MyClass.MyMethod", arg)
        ...
        return ...
        }
    ...
    }
```