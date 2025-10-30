### Building a List of Unique Values

One way is to use [object.AddUnique](<../../Language/Reference/Object/object.AddUnique.md>) which checks if the value is already in the list before adding it:

``` suneido
list.AddUnique(value)
```

However, for long lists the linear search to check gets slow.

Another option is to make the values the members (keys) and then use [object.Members](<../../Language/Reference/Object/object.Members.md>) to get a list. Named members are stored in a hash table so this is fast.

``` suneido
list[value] = true
...
list = list.Members()
```

However, this uses slightly more memory than a simple list.

Another option, especially if you need to sort the list anyway, is to use object.SortUnique! which sorts the list with object.Sort! and then removes duplicates in a single pass.

``` suneido
list.Add(value)
...
list.SortUnique!()
```

However, this may use more memory if there are a lot of duplicates since they won't be removed till the end.

Recommendations:

-	for a small number of values - AddUnique is the simplest
-	for a medium number of values or lots of duplicates - use members
-	for a large number of values with few duplicates (and if sorting is ok) - use SortUnique!