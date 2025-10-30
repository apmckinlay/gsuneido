### RetryBool

``` suneido
(maxRetries, min, block) => result
```

Retries the block up to maxRetries times with an exponential random fallback (using [RetrySleep](<RetrySleep.md>)).

The block must return true when it succeeds.

**Warning**: Beware of side effects from the block. e.g. If the block sends an email, multiple emails could end up getting sent.

After maxRetries it will return the result from the last call to the block

See also:
[Retry](<../../Database/Reference/Retry.md>)