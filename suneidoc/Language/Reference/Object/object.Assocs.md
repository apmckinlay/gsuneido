<div style="float:right"><span class="builtin">Builtin</span></div>

#### object.Assocs

``` suneido
() => sequence
```

Returns a sequence of objects each containing a member name and a value.

For example:

``` suneido
ob = #(12, 34, a: 56, b: 78)
ob.Assocs()
    => #(#(0, 12), #(1, 34), #(#b, 78), #(#a, 56))
```

**Note:** The named members are in no particular order since objects are implemented as hash tables.

The returned sequence is initially "virtual", i.e. it is simply an iterator over the original data. If you only iterate over the sequence (e.g. using a for-in loop) then an object containing all the values is never instantiated. However, if you access the sequence in most other ways (e.g. call a built-in object method on it) then a list will be created with the values from the sequence and that list is what will be used for any future operations. See also: [Basic Data Types > Sequence](<../../Basic Data Types/Sequence.md>)

Note: When testing on the WorkSpace, the variable display will trigger instantiation of any sequences in local variables.

**Warning**: Since the sequence just iterates over the original object, you can get "object modified during iteration" errors. This can be avoided by forcing a copy of the sequence with .Copy()


See also:
[object.Members](<object.Members.md>),
[object.Values](<object.Values.md>)
