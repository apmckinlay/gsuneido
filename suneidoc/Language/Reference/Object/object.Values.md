<div style="float:right"><span class="builtin">Builtin</span></div>

#### object.Values

``` suneido
() or (list:) or (named:) => sequence
```

Returns a sequence of the values from the object. If **list:** is specified then only the initial un-named members are included.

For example:

``` suneido
ob = #(12, 34, a: 56, b: 78)
ob.Values()
    => #(12, 34, 78, 56)
ob.Values(list:)
    => #(12, 34)
ob.Values(named:)
    => #(78, 56)
```

**Note:** The list: members are the consecutive integer members starting at zero. If a number is "missing" then subsequent numbers will go into the named: members. For example:

``` suneido
ob = Object(11, 22, 33, 44)
ob.Values(list:)
    => #(11, 22, 33, 44)
ob.Erase(2)
    => #(11, 22, 3: 44)
ob.Values(list:)
    => #(11, 22)
```

**Note:** The named members are in no particular order since objects are implemented as hash tables.

The returned sequence is initially "virtual", i.e. it is simply an iterator over the original data. If you only iterate over the sequence (e.g. using a for-in loop) then an object containing all the values is never instantiated. However, if you access the sequence in most other ways (e.g. call a built-in object method on it) then a list will be created with the values from the sequence and that list is what will be used for any future operations. See also: [Basic Data Types > Sequence](<../../Basic Data Types/Sequence.md>)

Note: When testing on the WorkSpace, the variable display will trigger instantiation of any sequences in local variables.

**Warning**: Since the sequence just iterates over the original object, you can get "object modified during iteration" errors. This can be avoided by forcing a copy of the sequence with .Copy()


See also:
[object.Members](<object.Members.md>),
[object.Assocs](<object.Assocs.md>)
