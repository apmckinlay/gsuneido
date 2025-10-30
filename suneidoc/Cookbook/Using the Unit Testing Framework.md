## Using the Unit Testing Framework

**Category:** Testing

**Problem**

You need to create automated *unit tests* for your code.

**Ingredients**

[Test](<../Language/Reference/Test.md>), TestRunner, TestRunnerGui

**Recipe**

In Suneido a unit test is a class in a library. The name must start with a capital letter (like all global library definitions) and end with "Test". The class must be derived (directly or indirectly) from the stdlib **Test** class. For example, the minimum (empty) test would be:

**`MyTest`**
``` suneido
Test
    {    }
```

You can run the test a variety of ways. One easy way is from Library View - with the cursor in your test (but nothing selected) hit F9 (or choose Run from the toolbar or menu). You should see an alert saying "ALL TESTS SUCCEEDED" or "FAILURES:" followed by a list of the methods that failed.

Alternately, you can "call" the test from the WorkSpace:

``` suneido
MyTest()
```

The most versatile method of running tests is to use the TestRunnerGui.
This is available from the WorkSpace on the Tools menu or via the "T"
toolbar button.  In the Test Runner you select the library of interest
(or all libraries) and it displays a list of the tests. You can then
choose to run them all, or just select a particular test to run.
Perhaps the best feature is that you can choose to debug a selected
test. This will bring up the debugger when an error occurs so you can
see exactly where it's happening, and examine variable values at that
point.

Of course, a useful test should actually test something! Test methods are
any methods whose names start with "Test". All such methods will be
run, in no particular order. In other words, don't write test methods
that depend on running in a particular order. Unlike some testing
frameworks, Suneido does not use special assert methods. You use the
standard [Assert](<../Language/Reference/Assert.md>) etc.. Exceptions are caught and recorded by
the testing framework. This means that test methods will only execute
up to the first exception. Here is a very simple test:

``` suneido
Test
    {
    Test_add()
        {
        Assert(1 + 1 is: 2)
        }
    Test_subtract()
        {
        Assert(1 - 1 is: 0)
        }
    }
```

This should, obviously, succeed. Try changing the 2 to a 3 to see it fail.

Test classes can also include Setup and Teardown methods. Setup is run once,
before any of the test methods. Teardown is called after all the test
methods.  (This is different from testing frameworks that call Setup
and Teardown before and after **each** test method.) If your test
creates any external artifacts, e.g. database records or tables, you
should clean these up in a Teardown method, rather than at the end of a
test method, since teardown will be called regardless whether the tests
succeed or fail. Here is a simple test using Setup and Teardown:

``` suneido
Test
    {
    Setup()
        {
        Database("create mytest (name, age) key(name)")
        }
    Test_add()
        {
        record = Record(name: "fred", age: 23)
        QueryOutput("mytest", record)
        Assert(QueryFirst("mytest") is: record)
        }
    Teardown()
        {
        try Database("destroy mytest")
        }
    }
```

Our convention is to make Setup the first method and Teardown the last. We
also often name test tables or files using the name of the test to help
ensure they're unique and to make it easier to trace where a stray
table or file came from.

The testing framework will complain if your test does not leave the
database with the same number of tables and records as when it started.

**Discussion**

XP (Extreme Programming, not Windows) popularized the idea of *test driven* development. In XP you write automated unit tests for all your code. The idea is that you only write code to fix *broken* tests. So before you create a new feature, you write a test for it. This forces you to think first about how the code will be *used*, before you get into the details of *how*
it will work. Similarly, when you find a bug in your program the first
thing you should do (before fixed it) is to write a test that catches
the bug. As well as ensuring that the bug will never come back without
getting caught, this also forces you to pin down exactly what situation
exposed the bug.

Suneido itself has a suite of unit tests. Some of these are built into the
executable and can be run from the command line via "suneido -tests".
Others are in stdlib and can be run as described above. In our office
we also have a scheduled task that runs every hour that grabs the
latest changes from our version control, runs the tests, and emails the
results to each of the programmers.  We also run the tests nightly on
our customer installations, again emailing the results back to
ourselves.

**See Also**

The *booklet* by Stefan Schmiedl which includes a description of how to use Suneido's testing framework.

SuneidoUnitLib,
an earlier testing framework by Randy Coulman, based more closely Kent
Beck's Extreme Programming xUnit testing framework, is available in the
library section. 

Other unit testing frameworks at [http://xprogramming.com/software.htm](<http://xprogramming.com/software.htm>)

Test Driven Development, Kent Beck, Addison-Wesley 2003

Testing Extreme Programming, Lisa Crispin & Tip House, Addison-Wesley 2003