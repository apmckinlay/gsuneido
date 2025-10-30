## A MockObject for Testing

**Category:** Coding

**Ingredients**

-	the special 
	[Default](<../Language/Classes/Default.md>) method
-	the testing framework


**Problem**

You want to test that some code does a certain sequence of calls on an object.

**Recipe**

Replace the object with a "mock" object for the test. A "mock" object simply checks the calls against what it is told to expect. Here is a simple Suneido mock object.

*`MockObject`*
``` suneido
class
    {
    New(calls)
        {
        .calls = calls
        .i = 0
        }
    Default(@args)
        {
        if not .calls.Member?(.i)
            throw "No call expected\n" $
                "got: " $ Display(args)
        if .calls[.i] isnt args
            throw "Incorrect call:\n" $
                "expected: " $ Display(.calls[.i]) $ "\n" $
                "got: " $ Display(args)
        ++.i
        }
    }
```

The special "Default" method is used to intercept all method calls on the object. The first argument to Default is the name of the method called.

**Example**

For example, if we wanted to test this:

*`MyFunc`*
``` suneido
function (ob)
    {
    ob.Init()
    ob.Heading("My Numbers")
    ob.Print(123, 456)
    }
```

We could use MockObject like this:

*`MyFunc_Test`*
``` suneido
Test
    {
    Test_main()
        {
        mock = MockObject(#(
            (Init)
            (Heading, "My Numbers")
            (Print, 123, 456)))
        MyFunc(mock)
        }
    }
```

Notice that we don't need any assertions - all the checking is done by MockObject. Try changing the expected calls supplied to MockObject and see how the error is caught.

**Going Further**

In some cases your code requires return values from the method calls. MockObject could easily be extended so that you supplied the return values for each call.

**See Also**

[MockObject](<../Language/Reference/MockObject.md>),
[FakeObject](<../Language/Reference/FakeObject.md>)

[http://www.mockobjects.com](<http://www.mockobjects.com>)   

Test Driven Development, Kent Beck   

Test Driven Development, David Astels   

Working Effectively With Legacy Code, Michael Feathers   
[Using the Unit Testing Framework](<Using the Unit Testing Framework.md>)

**Acknowledgments**

Thanks to Jeff Ferguson and Kim Dong for writing the initial version of [MockObject](<../Language/Reference/MockObject.md>).