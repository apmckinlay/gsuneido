<div style="float:right"><span class="toplinks"><a href="/suneidoc/Language/Reference/Spy/Methods">Methods</a></span></div>

### Spy

``` suneido
(target) => Spy
```

"target" can be a function, a method or a global name string.
   

**IMPORTANT:** Only use Spy in tests.

**NOTE:** You need to call spy.Close() to clean up a spy. Or you can call SpyManager().RemoveAll() to clean up all the spies.

**TIPS:** In a Test, the **.SpyOn** method creates a spy instance and also ensures it will clean up the spy at the end of the current test function scope automatically whether the test succeeds or fails. **Use this method instead of calling the class directly**.

#### Verifying Behavior with spy.CallLogs

``` suneido
spy = .SpyOn(Func)
Func(123, "xyz")
callLog = spy.CallLogs()
Assert(callLog isSize: 1)
Assert(callLog[0] is: Object(a: 123, b: "xyz"))
```

Where "Func" is any function or method with parameters "a" and "b".

#### Stubbing with spy.Return and spy.Throw

By default, Spied functions will call through when being invoked. You can specify a call to do the following things using **spy.Return/Throw**

return a value:

``` suneido
spy.Return("xyz")
```

return multiple values for consecutive calls:

``` suneido
spy.Return("first return", "second return", "third return")
```

return different values conditionally:

``` suneido
spy.Return("value 1", when: /* condition 1 block */)
spy.Return("value 2", when: /* condition 2 block */)
```

throw an exception:

``` suneido
spy.Throw("must be positive")
```

#### Testing a Method that Calls Other Methods and Global Functions

For example, you want to test this:

``` suneido
class
    {
    ...
    Method1()
        {
        if .Method2()
            Func1()
        }
    ...
    }
```

You can create two spies on .Method2 and Func1 and stub their return values separately

``` suneido
spy1 = .SpyOn(MyClass.Method2)
spy2 = .SpyOn(Func1)
spy1.Return(true)
MyClass.Method1()
Assert(spy2.CallLogs() isSize: 1)

spy1 = .SpyOn(MyClass.Method2)
spy2 = .SpyOn(Func1)
spy1.Return(false)
MyClass.Method1()
Assert(spy2.CallLogs() isSize: 0)
```