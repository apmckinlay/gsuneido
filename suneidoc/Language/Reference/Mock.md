<div style="float:right"><span class="toplinks"><a href="/suneidoc/Language/Reference/Mock/Methods">Methods</a></span></div>

### Mock

``` suneido
(cls = false) => mock
```

Creates a "mock" object that will accept any method calls. Used for tests.

``` suneido
mock = Mock()
mock.AnyMethod()
```

#### Verifying Behaviour with mock.Verify

``` suneido
mock = Mock()
mock.Func(123, "xyz")
mock.Verify.Func(123, "xyz")
```

Where "Func" is any method name.

If you verify a call and it was not invoked, then an exception will be thrown (and presumably the test will fail).

You can also verify that a call was <u>not</u> made with `mock.Verify.Never().Func(...)` (see below)

#### Stubbing with mock.When

By default, calling a method on a mock has no return value. You can specify a call to do the following things using **mock.When**

return a value:

``` suneido
mock.When.Func(123).Return("xyz")
```

call through when being invoked:

``` suneido
mock.When.Func(-123).CallThrough()
```

excute a block:

``` suneido
mock.When.Func(-123).Do({ /* block */ })
```

throw an exception:

``` suneido
mock.When.Func(-123).Throw("must be positive")
```

If the class name is passed into the Mock constructor, you can stub private methods without prefixing them with the class name.
<pre>
mock = Mock(<b>MyClass</b>)
mock.When.<b>func</b>().Return(123) // instead of mock.When.MyClass_func().Return(123)
</pre>

You can also specify multiple return values for consecutive calls.

``` suneido
mock = Mock()
mock.When.NextNumber().Return(0, 1, 2)
mock.NextNumber()
    => 0
mock.NextNumber()
    => 1
mock.NextNumber()
    => 2
mock.NextNumber()
    => 2 // additional calls return last value
```

#### Argument Matching

When you say `mock.Verify.Func(123, "xyz")` or `mock.When.Func(123, "xyz")` it will only apply when Func is called with those exact arguments.

**NOTE:** Mutable arguments are stored in Mock by references. Arguments passed are compared using the equals (is) method by default

If you want to match any number of arguments (including none) with any values, you can use `mock.Verify.Func([anyArgs:])` or `mock.When.Func([anyArgs:])`

To specify that one of the arguments may have any value, use `[any:]` as the argument.

It is also possible to use any of the Hamcrest style matchers used by [Assert](<Assert.md>). For readability, anyNumber, anyString, and anyObject are available that simply inherit from isNumber, isString, and isObject.

``` suneido
mock.When.Method('value1', [like: 'value']).CallThrough()

mock.Verify.Method('value1', [startsWith: 'value'])
```

#### Verifying Number of Calls

By default, `mock.Verify.Func()` checks that the call was made exactly once.

You can also use:
`mock.Verify.Never().Func()`
: Verify that the call was *not* made.

`mock.Verify.Times(n).Func()`
: Verify that the call was made exactly n times.

`mock.Verify.AtLeast(n).Func()`
: Verify that the call was made n or *more* times.

`mock.Verify.AtMost(n).Func()`
: Verify that the call was made n or *less* times.

If the class is passed into the Mock constructor, you can verify private methods without prefixing them with the class name.

#### Testing a Method that Calls Other Methods

For example, you want to test this:

``` suneido
class
    {
    ...
    Method1()
        {
        if .Method2()
            .Method3()
        }
    ...
    }
```

You can use [object.Eval](<Object/object.Eval.md>) to call the method to be tested in the context of the mock.

``` suneido
mock = Mock()
mock.When.Method2().Return(true)
mock.Eval(MyClass.Method1)
mock.Verify.Method3()

mock = Mock()
mock.When.Method2().Return(false)
mock.Eval(MyClass.Method1)
mock.Verify.Never().Method3()
```

Or, you can stub CallThrough() on a method and then call it in the context of the mock, the class needs to be passed into the Mock constructor

``` suneido
mock = Mock(MyClass)
mock.When.Method2().Return(true)
mock.When.Method1().CallThrough()
MyClass.Method1()
mock.Verify.Method3()
```

If you pass the class when constructing the mock, then you can reference private methods directly:
<pre>
mock = Mock(<b>MyClass</b>)
mock.When.<b>method2</b>().Return(false)
mock.When.<b>method1</b>().CallThrough()
mock.<b>method1</b>()
mock.Verify.method3()
</pre>

Otherwise, you will need to prefix them with the class name. For example:

``` suneido
mock = Mock()
mock.When.MyClass_method2().Return(false)
mock.Eval(MyClass.MyClass_method1)
```

**NOTE:** If the class is passed into Mock object, all value members are copied to the mock object automatically.

#### See Also

[MockObject](<MockObject.md>),
[FakeObject](<FakeObject.md>),
[From the Couch 12 - Mockito for Suneido](<https://suneido.com/from-the-couch-12-mockito-for-suneido/>) (on the web site)