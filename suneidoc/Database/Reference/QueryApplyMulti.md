### QueryApplyMulti

``` suneido
(query, block, update)
(query, update) { |x| ... }
```

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

QueryApplyMulti is not necessary for read-only i.e. update should normally be true. The reason for the update argument is so you can use the same code for both read-only and update purposes.