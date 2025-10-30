#### object.Swap

``` suneido
(i, j) => this
```

Swaps the value at **i** with the value at **j**. Can be used with numeric or named members.

For example:

``` suneido
Object(12, 34, a: 56, b: 78).Swap(1, 'a')
    => #(12, 56, a: 34, b: 78)
```

**Note:** Swap modifies the object it is applied to, 
it does not create a new object.