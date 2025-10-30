<div style="float:right"><span class="builtin">Builtin</span></div>

#### object.Add

``` suneido
(value ...) => this
(value ... at: position) => this
```

Add values to an object. If no position is specified, the values are added to the end of the un-named members. For example:

``` suneido
Object(12, 34).Add(56, 78) => Object(12, 34, 56, 78)
```

If the position is *within *the un-named members, the values are *inserted *there.

``` suneido
Object(12, 78).Add(34, 56, at: 1) => Object(12, 34, 56, 78)
```

If the position is an integer, but is not within the un-named members, the values will be added at consecutive positions.  Note: it is possible for this to overwrite existing members.

``` suneido
Object().Add(12, 34, 56, at: 10) => Object(10: 12, 11: 34, 12: 56)
```

If the position is a member name, then only a single value may be specified and Add simply sets that member to the value.

``` suneido
Object().Add(123, at: 'num') => Object(num: 123)
```

**Note:** Add modifies the object it is applied to, it does not create a new object.

**Note:** Add is optimized for x.Add(@y) where x and y have no named members.

See also: [object.AddTo](<object.AddTo.md>)