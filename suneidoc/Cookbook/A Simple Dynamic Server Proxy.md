## A Simple Dynamic Server Proxy

**Category:** Coding

**Ingredients**

-	the special Default method
-	the built-in Display function
-	string.ServerEval


**Problem**

You want to interface with a remote object or class without having to write a custom "proxy" class.

**Recipe**

We can use the special Default method to redirect calls to a remote object using ServerEval.

*`ServerEvalProxy`*
``` suneido
class
    {
    New(remote)
        {
        .remote = remote
        }
    Default(@args)
        {
        call = .remote $ '.' $ args[0] $ Display(args[1..])[1..]
        return call.ServerEval()
        }
    }
```

The "remote" argument is the class or object on the server.

Display is used to convert the arguments to a string. Slice is used to skip the first argument (the name of the method called). Substr(1) is used to skip the leading '#' that Display will include.

A nice thing about ServerEval is that it will work even when you're running standalone, without a server. In this case it will be equivalent to Eval.

Note: Using Display with ServerEval will only handle arguments that can safely be converted to strings and back. Fortunately, this includes the most common types such as numbers, strings, dates, and objects containing such types.

You can create a local proxy with:

``` suneido
proxy = ServerEvalProxy("MyServerObject")
```

If you want a global proxy, you can derive a class, see the example below.

**Unit Test**

Here is a unit test for ServerEvalProxy:

*`ServerEvalProxy_Test`*
``` suneido
Test
    {
    Test_main()
        {
        proxy = ServerEvalProxy("ServerEvalProxy_Test")
        Assert(proxy.A() is: 123)
        Assert(proxy.B(20, "th century ", #20000101) is: "20th century 2000-1-1")
        }
    A()
        { return 123 }
    B(number, string, date)
        { return number $ string $ date.Format('yyyy-M-d') }
    }
```

This uses the "self shunt" testing pattern where the testing class itself is used as an object in the test. Since we need a global class,  using self shunt saves us from having to define a separate library record.

**Example**

I developed this code for use with a lock manager for pessimistic locking in AccessControl. AccessControl calls it like:

``` suneido
LockManager.Lock(key)
```

LockManager is defined as:

``` suneido
ServerEvalProxy
    {
    New()
        { super("LockManagerImpl()") }
    }
```

LockManagerImpl contains the actual lock manager implementation and runs on the server (when running client-server). It uses CallClass to implement the Singleton pattern:

``` suneido
CallClass()
    {
    if not Suneido.Member?("LockManager")
        Suneido.LockManager = new this
    return Suneido.LockManager
    }
```

**See Also**

Design Patterns by Gamma, Helm, Johnson, Vlissides includes the Proxy and Singleton patterns.

Head First Design Patterns by Freeman & Freeman covers these patterns in a more "friendly" style.

Test Driven Development by Kent Beck includes the Self Shunt testing pattern.