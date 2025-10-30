### DoWithTran

``` suneido
(t, block, update = false)
```

If no transaction is passed in, a new transaction will be created.  The transaction is then passed to the block.

For example:

``` suneido
function (t = false)
    {
    DoWithTran(t, update:)
        { |t| 
        x = t.Query1('mytable where key = 1')
        x.date = Date()
        x.Update()
        }
    }
```

See also:
[record.DoWithTran](<Record/record.DoWithTran.md>)