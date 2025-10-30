### Methods

When you retrieve a method from an object or class, you get a "Method" value that is bound to the object that it was retrieved from. When you call this value, it is as if you had called that method on that object or class. i.e. it has access to **this** and member variables as it normally would.

For example:

``` suneido
c = class
    {
    member: 123
    Fn() { return .member }
    }

c.Fn()
    => 123

method = c.Fn
method()
    => 123
```

[object.Eval](<../Reference/Object/object.Eval.md>) can be used to override this behavior and call the method for a different instance or class.