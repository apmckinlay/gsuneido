#### string.Ellipsis

``` suneido
(maxLength, atEnd = false) => string
```

If the string size is less than or equal to maxLength, it is returned as is.

Otherwise, the excess is replaced with "..." ether in the middle or the end.

For example:

``` suneido
"hello world".Ellipsis(4)
	=> "he...ld"

"hello world".Ellipsis(5, atEnd:)
	=> "hello..."
```