### WithQuery

``` suneido
(@args)
```

Does:

``` suneido
Transaction(read:)
    { |t|
    return t.Query(@args)
    }
```

For example:

``` suneido
WithQuery('tables')
    { |q|
    first = q.Next()
    second = q.Next()
    }
```

See also:
[transaction.Query](<Transaction/transaction.Query.md>)