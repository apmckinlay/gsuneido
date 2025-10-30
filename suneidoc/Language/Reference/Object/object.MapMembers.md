#### object.MapMembers

``` suneido
(block) => object
```

Similar to [object.Map](<object.Map.md>) except that block is used to transform the keys rather than the values (i.e. it transforms `object.Members()`, not `object.Values()`).

block can be anything that can be called e.g. block, function, class, or instance.

For example:

``` suneido
#(a b c).MapMembers({ ("A".Asc() + it).Chr() })
    => #(A: "a", B: "b", C: "c")
```