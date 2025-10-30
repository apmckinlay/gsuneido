<div style="float:right"><span class="builtin">Builtin</span></div>

#### Suneido.StrictCompare

``` suneido
(boolean)
```

StrictCompare starts as false. If StrictCompare is set to true, then comparisons (\<, \<=, >, >=) between different types will throw an exception.

The reason for this is that comparisons between types are ambiguous. You might write > 0 for positive numbers, but it also unintentionally includes e.g. dates. If it's intended for positive numbers, x > 0 can be rewritten as: Number?(x) and x > 0 (or use [Positive?](<../Positive?.md>))

[By](<../By.md>), [Gt](<../Gt.md>), [Min](<../Min.md>), [object.Min](<../Object/object.Min.md>), [object.MinWith](<../Object/object.MinWith.md>), [Max](<../Max.md>), [object.Max](<../Object/object.Max.md>), [object.MaxWith](<../Object/object.MaxWith.md>), [object.Sort!](<../Object/object.Sort!.md>) are exempt from Suneido.StrictCompare

The goal is to eventually make StrictCompare always enabled.


See also:
[Cmp](<../Cmp.md>),
[Gt](<../Gt.md>),
[Min](<../Min.md>),
[Max](<../Max.md>),
[object.Min](<../Object/object.Min.md>),
[object.MinWith](<../Object/object.MinWith.md>),
[object.Max](<../Object/object.Max.md>),
[object.MaxWith](<../Object/object.MaxWith.md>)
