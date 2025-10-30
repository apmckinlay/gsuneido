#### object.GetInit

``` suneido
(member, value or block) => value
```

If the object does not contains **member**, it will be set to value. Then the value of member will be returned. For example:

``` suneido
#(a: 12, b: 34).GetInit("b", 99) => 34

x = [a: 12, b: 34]
x.GetInit("c", 99) => 99
x => [a: 12, b: 34, c: 99]
```

If the value is a block, then it will be evaluated, but only if needed to initialize member. For example:

``` suneido
x = Suneido.GetInit("mycache", { Query1(...) })
```

If you did Suneido.GetInit("mycache", Query1(...)) without a block it would do the query every time it was called, even though it would only use the result the first time.

If multiple threads try to initialize at the same time, the block may be evaluated more than once. But the member will still only be initialized once and all the GetInit calls will return the same value. e.g.

``` suneido
Thread({ x = Suneido.GetInit(#foo, Timestamp) })
Thread({ y = Suneido.GetInit(#foo, Timestamp) })
Thread({ z = Suneido.GetInit(#foo, Timestamp) })
```

Timestamp may be called more than once, but one of the threads will "win" and x, y, and z will all get the same Timestamp value.

See also: [object.GetDefault](<object.GetDefault.md>)