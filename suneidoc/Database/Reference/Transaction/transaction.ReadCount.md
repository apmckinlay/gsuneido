#### transaction.ReadCount

``` suneido
() => number
```

Returns the number of "reads" that have been performed in the transaction. This is useful for debugging issues where you are exceeding the limits and getting a "too many reads" exception.

A "read" is defined as accessing a range of an index. The limit is 10,000 reads on any one index.

ReadCount returns the total reads for all indexes, therefore it may report a number greater than 10,000.

See also: [transaction.WriteCount](<transaction.WriteCount.md>)