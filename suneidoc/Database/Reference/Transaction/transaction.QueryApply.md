#### transaction.QueryApply

See: [QueryApply](<../QueryApply.md>)

t.QueryApply takes an additional readonly: option

Normally QueryApply in an update transaction adds a sort by a key. This is to reduce problems when records are modified during iteration. However, in some cases, t.QueryApply is used within an update transaction, but that particular QueryApply is not updating the records being iterated over. In this case you can specify readonly: and QueryApply will not add a key sort. This allows better query optimization.

**Note**: the readonly: option is not enforced. If you specify readonly: but actually do update the records being iterated, then you will get unpredictable behavior.