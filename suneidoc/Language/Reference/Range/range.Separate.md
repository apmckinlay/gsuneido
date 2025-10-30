### range.Separate

``` suneido
( range ) => object
```

where **range** is a Range that is included in this.

Returns an object containing 2 ranges, where the first range is formed by the *low* values of the current range and the **range** passing in; and the second range is formed by the *high* values of the current range and the **range** passing in.

For example:

``` suneido
range = Range(2, 7)
range.Separate(Range(3, 5)) => #(Range(2, 3), Range(5, 7))
```