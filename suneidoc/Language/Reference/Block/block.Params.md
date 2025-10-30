<div style="float:right"><span class="builtin">Builtin</span></div>

#### block.Params

``` suneido
() => string
```

Returns a string containing the parameter list for the block.

Works for library functions, built-in functions, and blocks.

For example:

``` suneido
block = {|x, y| /*...*/ }
block.Params()
    => "(x,y)"
```

See also: [function.Params](<../Function/function.Params.md>)