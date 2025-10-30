### MockObject

``` suneido
(calls)
```

MockObject is used to check a sequence of calls (and arguments) on an object against a list of expected **calls** and arguments passed to the MockObject. An exception is thrown if any of the calls or arguments do not match the list provided. MockObject uses the [Default](<../Classes/Default.md>) method to keep track of the method calls on the object.

A call is specified by an object containing the name of the method followed by the arguments. For example:

``` suneido
#(Calc, 123, 456, round:)
```

A call specified like this will have no return value. If needed you can specify the return value:

``` suneido
#((Calc, 123, 456, round:) result: 1.7)
```

For example, if we wanted to test this:

``` suneido
MyFunc

function (ob)
    {
    ob.Init()
    ob.Heading("My Numbers")
    return ob.Print(123, 456)
    }
```

We could use MockObject like this:

``` suneido
MyFunc_Test

Test
    {
    Test_main()
        {
        mock = MockObject(#(
            (Init)
            (Heading, "My Numbers")
            ((Print, 123, 456)), result: "My Numbers 123 456"))
        MyFunc(mock)
        }
    }
```

Notice that we don't need any assertions - all the checking is done by MockObject.

For a different approach see [Mock](<Mock.md>)

See also: [FakeObject](<FakeObject.md>) and
[A Mock Object for Testing](<../../Cookbook/A Mock Object for Testing.md>)