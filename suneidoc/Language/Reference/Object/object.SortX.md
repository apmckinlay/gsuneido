<div style="float:right"><span class="builtin">Builtin</span></div>

#### object.Sort!
<pre>
(less_function = <a href="/suneidoc/Language/Reference/Lt">Lt</a>) => this
</pre>

Sorts the *list* members of the object (consecutive members starting at 0).  
It has no effect on named members.

For example:

``` suneido
Object(23, 76, 34, 98, 56).Sort!() => #(23, 34, 56, 76, 98)
```

A comparison function can be supplied. This is normally required when the values being sorted are objects or records. The comparison function is called with two values and should return true if the first value is less than the second (not if they are equal) and false otherwise. The comparison function can be anything callable e.g. a block, function, or method. For example:

``` suneido
list = Object(#(3, "Three"), #(1, "One"), #(2, "Two"))
list.Sort!({ |x,y| x[0] < y[0] })
    => Object(#(1, "One"), #(2, "Two"), #(3, "Three"))
```

To sort in reverse order, pass [Gt](<../Gt.md>) as the comparison function.

Sort! is a stable sort so it will maintain the order of *equal* items.

**Note:** Sort! modifies the object it is applied to, it does not create a new object.


See also:
[By](<../By.md>),
[object.Sorted?](<object.Sorted?.md>),
[object.SortWith!](<object.SortWith!.md>)
