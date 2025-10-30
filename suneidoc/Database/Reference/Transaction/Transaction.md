#### Transaction

``` suneido
(read:) => transaction
(update:) => transaction
(read:, block) => value
(update:, block) => value
```

Create a new database transaction. For example:

``` suneido
t = Transaction(read:)
x = t.Query('stdlib').Next()
t.Complete()
return x
```

If a block (or other callable value) is supplied, it is called with the transaction as its argument, the transaction is automatically Complete'd when it returns, and the return value from the block is returned. If the block throws an exception, Rollback is called on the transaction. This is a good way to ensure that transactions are completed. For example: 

``` suneido
x = Transaction(read:)
    {|t| t.Query('stdlib').Next() }
```

The block form throws an exception ("Transaction: block commit failed") if the transaction fails to complete.