<div style="float:right"><span class="builtin">Builtin</span></div>

#### object.Iter

``` suneido
() => object_iterator
```

Returns an iterator for the object that provides:

``` suneido
.Next() => member or iterator
.Iter() => iterator
.Rewind()
```

Iter is used by **`for`**.

``` suneido
for (m in ob)
```

is equivalent to:

``` suneido
iter = ob.Iter()
while (iter isnt m = iter.Next())
```

See also:
[string.Iter](<../String/string.Iter.md>)