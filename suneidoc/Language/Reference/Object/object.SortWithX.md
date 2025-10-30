#### object.SortWith!

``` suneido
(callable) => this
```

Does:

``` suneido
.Sort!({|x,y| fn(x) < fn(y) })
```

For example:

``` suneido
['a','B','c','D'].SortWith!(#Lower)

    => ["a", "B", "c", "D"]
```


See also:
[By](<../By.md>),
[object.Sort!](<object.Sort!.md>),
[object.Sorted?](<object.Sorted?.md>)
