### range.Empty?

``` suneido
( ) => boolean
```

Returns true if low and high are both false, else returns false.

For example:

``` suneido
Range().Empty?() => true
Range(2, 7).Empty?() => false
```