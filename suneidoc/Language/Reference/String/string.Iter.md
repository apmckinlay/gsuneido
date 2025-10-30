<div style="float:right"><span class="builtin">Builtin</span></div>

#### string.Iter

``` suneido
() => string_iterator
```

Returns an iterator for the string.  String iterators provide:

``` suneido
.Next() => next character or iterator itself if no more characters
```

String iterators work with `for` loops.

``` suneido
for (char in string)
```

is the same as:

``` suneido
for (iter = string.Iter(); iter isnt char = iter.Next(); )
```

See also: [object.Iter](<../Object/object.Iter.md>)