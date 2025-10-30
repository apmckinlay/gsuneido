#### object.Fold

``` suneido
(init, block)
(init) { ... }
```

Calls the block for each member of the object, passing it the accumulated value (starting with **init**). The block must return the updated accumulated value. This is a left fold. [object.Reduce](<object.Reduce.md>) is similar but it takes the first element as the initial value.

For example, to sum the values in an object: (but in practice, use [object.Sum](<object.Sum.md>))

``` suneido
#(123, 456, 789).Fold(0, {|sum x| sum += x })
    => 1368
```

Or to concatenate the values in an object: (but in practice, use [object.Join](<object.Join.md>))

``` suneido
#("abc", "def", "ghi").Fold("", {|s x| s $= x })
    => "abcdefghi"
```

Or to average the values: 

``` suneido
ob = #(10, 20, 60).Fold(Object(sum: 0, n: 0))
    { |ob x|
    ob.sum += x
    ++ob.n
    ob
    }
ob.sum / ob.n
    => 30
```


See also:
[Compose](<../Compose.md>),
[Curry](<../Curry.md>),
[Memoize](<../Memoize.md>),
[MemoizeSingle](<../MemoizeSingle.md>),
[object.Any?](<object.Any?.md>),
[object.Drop](<object.Drop.md>),
[object.Every?](<object.Every?.md>),
[object.Filter](<object.Filter.md>),
[object.FlatMap](<object.FlatMap.md>),
[object.Map](<object.Map.md>),
[object.Map!](<object.Map!.md>),
[object.Map2](<object.Map2.md>),
[object.Nth](<object.Nth.md>),
object.PrevSeq,
[object.Reduce](<object.Reduce.md>),
[object.Take](<object.Take.md>),
[object.Zip](<object.Zip.md>)
