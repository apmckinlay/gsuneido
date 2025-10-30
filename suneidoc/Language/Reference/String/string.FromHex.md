<div style="float:right"><span class="builtin">Builtin</span></div>

#### string.FromHex

``` suneido
() => string
```

Converts each pair of characters to a single byte.

Will throw an exception if not valid hex.

For example:

``` suneido
"616263".FromHex()
	=> "abc"
```

See also: [string.ToHex](<string.ToHex.md>)