#### transaction.WriteCount

``` suneido
() => number
```

Returns the number of "writes" that have been performed in the transaction. This is useful for debugging issues where you are exceeding the limits and getting a "too many writes" exception.

A "write" is either adding a record, updating a record, or removing a record. The limit is 10,000 writes per transaction.

See also: [transaction.ReadCount](<transaction.ReadCount.md>)