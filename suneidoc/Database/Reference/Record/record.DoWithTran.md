<div style="float:right"><span class="builtin">Builtin</span></div>

#### record.DoWithTran

``` suneido
() { |t| ... }
(update:) { |t| ... }
```

Rules should normally use their enclosing transaction, available via 
[record.Transaction](<record.Transaction.md>).
However, records do not always have an enclosing transaction (e.g. a new record).
This leads to the following type of code:

``` suneido
if (false is t = .Transaction())
    t = Transaction(read:)
...
if (.Transaction() is false)
    t.Complete()
```

DoWithTran lets you write this as simply:

``` suneido
.DoWithTran()
    { |t|
    ...
    }
```

or if you require an update transaction:

``` suneido
.DoWithTran(update:)
    { |t|
    ...
    }
```

**Note:** `DoWithTran(update:)` will not change an enclosing transaction from read-only to update.

See also:
[DoWithTran](<../DoWithTran.md>)