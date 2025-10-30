#### string.ReplaceSubstr

``` suneido
(start, length, replacement) => string
```

Returns a new string, with the substring starting at the specified position (zero based) 
and of the specified length replaced with the specified replacement argument.
The replacement string can be a different length than the substring it is replacing
(including zero length).

For example:

``` suneido
"hello there world".ReplaceSubstr(6, 5, "wonderful")
	=> "hello wonderful world"

"hello world".ReplaceSubstr(6, 0, "cruel ")
	=> "hello cruel world"
```

Note: The original string is not modified.

See also:
[string.Replace](<string.Replace.md>)