#### object.Flatten

``` suneido
() => object
```

Returns a new object, with each value copied,
recursively extracting the members of each value that is an object. 

For example:

``` suneido
#(123, "abc", name: "Fred", age: 25).Flatten() => #(123, "abc", "Fred", 25)
#(12, (34, 56), 78) => #(12, 34, 56, 78)
#(12, (34, (56, (78)))) => #(12, 34, 56, 78)
```

**Note:** Any member names are lost.