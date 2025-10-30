### QueryApply

``` suneido
(query, block, update = false, dir = 'Next')
(query, update = false, dir = 'Next') { |x| ... }
```

To process the records in reverse order, specify dir: 'Prev'

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

QueryApply(update:) will add a key sort to the query automatically, and will throw an exception if the query already has a sort but it's not a key.

See also: [QueryApply1](<QueryApply1.md>)