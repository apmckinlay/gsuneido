<div style="float:right"><span class="builtin">Builtin</span></div>

#### object.Unique!

``` suneido
() => this
```

Deletes *adjacent* duplicates from the *list* members of the object (consecutive members starting at 0).

Note: To eliminate *all* duplicates use [object.Sort!](<object.Sort!.md>) first.

For example:

``` suneido
Object(1, 2, 2, 3, 2).Unique!()
    => #(1, 2, 3, 2)

Object(2, 2, 2).Unique!()
    => #(2)
```

See also:
[object.UniqueValues](<object.UniqueValues.md>)