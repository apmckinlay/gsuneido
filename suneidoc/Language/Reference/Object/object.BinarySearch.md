<div style="float:right"><span class="builtin">Builtin</span></div>

#### object.BinarySearch

``` suneido
(value, block = less_than) => number
```

Returns the *first* position where the comparison returns false. i.e. the first position in a sorted list where the value could be inserted and maintain the sort order.

For example:

``` suneido
#(1, 2, 2, 3).BinarySearch(2)
    => 1
```

A comparison function can be supplied. This is normally required when the values being sorted are objects or records. The comparison function is called with two values and would normally return true if the first value is less than the second. The comparison function can be anything callable e.g. a block, function, or method.

BinarySearch with a comparison of less-than is equivalent to lower bound. With a comparison of less-than-or-equal-to it is equivalent to upper bound.

For example:

``` suneido
#(1, 2, 2, 3).BinarySearch(2, {|x,y| x <= y})
    => 3
```

**Note**: The object must be ordered with respect to the comparison function e.g. with [object.Sort!](<object.Sort!.md>)

See also: [object.BinarySearch](<object.BinarySearch.md>)