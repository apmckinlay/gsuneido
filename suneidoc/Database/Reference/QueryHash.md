### QueryHash

``` suneido
(query, details = false) => result
```

Returns a checksum of the columns and data resulting from a query. This is useful to confirm the results of the query are the same. The checksum is not affected by the order of the columns or data (since that can vary depending on the strategy.

If details is false the result is a single number. If details is true the result will be a string giving several values, like:

``` suneido
nrows 9 hash 1839799316 ncols 12 hash 3747039310
```