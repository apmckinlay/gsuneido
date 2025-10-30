### QueryApply1

``` suneido
(query, block)
(query) { |x| ... }
```

Calls the block for one record from the query in a standalone update transaction.

For example:

``` suneido
QueryApply1("mytable where mykey is 3")
    { |x| 
    x.name = x.name.CapitalizeWords()
    x.Update()
    }
```

This is equivalent to:

``` suneido
t = Transaction(update:)
x = t.Query1("mytable")
x.name = x.name.CapitalizeWords()
x.Update()
t.Complete()
```

See also: [QueryApply](<QueryApply.md>)