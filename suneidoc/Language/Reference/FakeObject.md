### FakeObject

``` suneido
(methods)
```

FakeObject is useful for testing. You can specify methods that should be handled and what they should do or return.

If a method is specified as simply a value, it will ignore any arguments and always return that value. For example:

``` suneido
ob = FakeObject(Calc: 0)
ob.Calc()
    => 0
ob.Calc(123, round: 2)
    => 0
```

A call can also be specified as a function. For example:

``` suneido
ob = FakeObject(FieldToPrompt: function (field) { return field.Capitalize() })
ob.FieldToPrompt('name')
    => "Name"
```

In this case, the arguments must be valid for the function parameters. You can use function(@args) to accept any arguments.

For example, if you wanted to test a function (or method) that needed a transaction argument and called QueryDo and Complete, you could make a "fake" transaction with:

``` suneido
faketran = FakeObject(QueryDo: true, Complete: true)
```

See also: [MockObject](<MockObject.md>)