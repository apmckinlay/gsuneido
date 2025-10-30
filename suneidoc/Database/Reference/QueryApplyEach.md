### QueryApplyEach

``` suneido
(query, block)
(query) { |x| ... }
```

Calls the block for each record from the query in an update transaction. It uses a read-only QueryApply to iterate over the query and a separate update transaction to process each record. This is useful when the process for each record is too large for QueryApplyMulti.

One of three related functions to iterate through a query and process each record.

Normal usage is like:

``` suneido
QueryApply/Multi/Each(query)
    { |x|
    ...
    }
```
[QueryApply](<QueryApply.md>)
: Use for read-only or for updating less than 10,000 records. Uses a single update transaction so it is atomic.

[QueryApplyMulti](<QueryApplyMulti.md>)
: Use for updating more than 10,000 records. Uses a separate update transaction for every 100 records.

[QueryApplyEach](<QueryApplyEach.md>)
: Use when there is a large amount of processing for each record. The query is iterated in a read-only transaction and each record is processed in its own RetryTransaction.