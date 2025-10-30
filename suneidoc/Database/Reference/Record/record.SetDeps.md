#### record.SetDeps

``` suneido
(member, deps)
```

Sets the dependencies for the specified member (field) to the comma separated list.

**Note**: SetDeps adds to any existing dependencies.

For example, if you had a Rule_sum:

``` suneido
function ()
    {
    return .a + .b
    }
```

Then:

``` suneido
r = Record(a: 2, b: 3, sum: 5)
r.a = 4
r.sum
    => 5
```

Since sum was set manually, there were no dependencies, and so modifying a did not cause sum to be updated.

But with SetDeps, it works properly:

``` suneido
r = Record(a: 2, b: 3, sum: 5)
r.SetDeps("sum", "a,b")
r.a = 4
r.sum
    => 7
```

See also: [record.GetDeps](<record.GetDeps.md>)