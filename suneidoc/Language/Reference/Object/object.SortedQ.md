#### object.Sorted?
<pre>
(less_function = <a href="/suneidoc/Language/Reference/Lt">Lt</a>) => true or false
</pre>

Returns true if the the *list* members of the object (consecutive members starting at 0) are already in sorted order, false otherwise.

Does not modify the object.

For example:

``` suneido
#(12, 34, 56).Sorted?()
    => true

#((name: "Fred", age: 35), (name: Sue, age: 25)).Sorted?(By(#age))
    => false
```


See also:
[By](<../By.md>),
[object.Sort!](<object.Sort!.md>),
[object.SortWith!](<object.SortWith!.md>)
