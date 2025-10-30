### Retry

``` suneido
(block, maxRetries = 10, minDelayMs = 2)
```

Retries the block up to maxRetries times with an exponential random fallback (using [RetrySleep](<RetrySleep.md>)).

The block must throw if it fails. (Any return value is ignored.)

**Warning**: Beware of side effects from the block. e.g. If the block sends an email, multiple emails could end up getting sent.

After maxRetries it will throw "Retry failed - too many retries, last error: ...". (The specific errors from earlier retries are not tracked.)

See also:
[RetryTransaction](<../../Database/Reference/RetryTransaction.md>)