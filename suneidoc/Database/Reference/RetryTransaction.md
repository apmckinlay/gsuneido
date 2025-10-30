### RetryTransaction

``` suneido
(block)
() { |t| ... }
```

Executes the block in an update transaction. If the transaction fails to commit, the block is executed again (in a new transaction). After 10 tries it throws "RetryTransaction: too many retries".

**Warning**: Beware of side effects from the block. e.g. If the block sends an email, multiple emails could end up getting sent.

For example:

``` suneido
RetryTransaction()
    { |t|
    x = t.Query1('nextnumber')
    n = x.number++
    x.Update()
    }
```

See also:
[Retry](<../../Language/Reference/Retry.md>)