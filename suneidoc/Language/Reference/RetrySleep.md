### RetrySleep

``` suneido
(i, min = 2)
```

Sleeps for a random amount of time in the range min \<< i to 2 * min \<< i.

i is normally the retry count (starting at zero).

Used by [Retry](<Retry.md>) and [RetryTransaction](<../../Database/Reference/RetryTransaction.md>)