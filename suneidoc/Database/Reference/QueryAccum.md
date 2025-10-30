<div style="float:right"><span class="deprecated">Deprecated</span></div>

### QueryAccum

``` suneido
(query, init, block)
(query, init) { ... }
```

QueryAccum calls the block for each record from the query, passing it the accumulated value (starting with **init**) and the record. The block must return the updated accumulated value.

For example, to sum the values in an object:

``` suneido
QueryAccum('tables', 0, { |sum rec| sum += rec.totalsize })
    => 4660686
```

Or to average the values:

``` suneido
ob = QueryAccum("tables", Object(sum: 0, n: 0))
    { |ob rec|
    ob.sum += rec.totalsize
    ++ob.n
    ob
    }
ob.sum / ob.n
    => 34762.48148148
```

**Note**: The "value" of a block is the value of its last statement. On the other hand, an actual "return" will return from the function containing the block.

QueryAccum is similar to Smalltalk's inject:into:, C++ STL accumulate, or Lisp's reduce.

See also:
[object.Accum](<../../Language/Reference/Object/object.Accum.md>)