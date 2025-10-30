#### record.GetDeps

``` suneido
(member) => string
```

Returns a string containing a comma separated list of the dependencies for a given member ("" if there are none).

For example, if you had a Rule_sum:

``` suneido
function ()
    {
    return .a + .b
    }
```

Then you would get:

``` suneido
r = Record(a: 123, b: 456)
r.sum
r.GetDeps("sum")
    => "a,b"
```

record.GetDeps returns the same list that would be saved if the dependencies were stored in an "_deps" field.

See also: [record.SetDeps](<record.SetDeps.md>)