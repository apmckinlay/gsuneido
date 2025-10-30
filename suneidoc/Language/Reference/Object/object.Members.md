<div style="float:right"><span class="builtin">Builtin</span></div>

#### object.Members

``` suneido
() or (list:) or (named:) or (all:) => sequence
```

Returns a sequence of the member names of the object. The initial un-named  members are listed first, in order, then the named members. If **list:** is specified then only the initial un-named members are included.

For classes or instances, you can do Members(all:) to include inherited members.

For example:

``` suneido
ob = #(12, 34, a: 56, b: 78)
ob.Members()
    => #(0, 1, #b, #a)
ob.Members(list:)
    => #(0, 1)
ob.Members(named:)
    => #(#b, #a)
```

**Note:** The list: members are the consecutive integer members starting at zero. If a number is "missing" then subsequent numbers will go into the named: members. For example:

``` suneido
ob = Object(11, 22, 33, 44)
ob.Members(list:)
    => #(0, 1, 2, 3)
ob.Erase(2)
    => #(11, 22, 3: 44)
ob.Members(list:)
    => #(0, 1)
```

**Note:** The named members are in no particular order since objects are implemented as hash tables.

The returned sequence is initially "virtual", i.e. it is simply an iterator over the original data. If you only iterate over the sequence (e.g. using a for-in loop) then an object containing all the values is never instantiated. However, if you access the sequence in most other ways (e.g. call a built-in object method on it) then a list will be created with the values from the sequence and that list is what will be used for any future operations. See also: [Basic Data Types > Sequence](<../../Basic Data Types/Sequence.md>)

Note: When testing on the WorkSpace, the variable display will trigger instantiation of any sequences in local variables.

**Warning**: Since the sequence just iterates over the original object, you can get "object modified during iteration" errors. This can be avoided by forcing a copy of the sequence with .Copy()


See also:
[object.Values](<object.Values.md>),
[object.Assocs](<object.Assocs.md>)
