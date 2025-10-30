#### record.SafeMembers

``` suneido
() => object
```

`rec.SafeMembers()` is equivalent to `rec.Members().Copy()`

It is necessary to avoid "object modified during iteration" errors due to members being created by rules.

For convenience, SafeMembers is also defined on objects, in which case it is equivalent to .Members (Objects don't have rules so they don't have this issue.)

See also: [object.Members](<../../../Language/Reference/Object/object.Members.md>)