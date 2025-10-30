### QueryDo

``` suneido
(request [, field: value ...])
```

Executes the insert, update, and delete request using [transaction.QueryDo](<Transaction/transaction.QueryDo.md>) in a standalone [RetryTransaction](<RetryTransaction.md>).

For example:

``` suneido
QueryDo("delete customers where sales < 100")
```