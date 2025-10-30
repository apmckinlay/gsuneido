<div style="float:right"><span class="builtin">Builtin</span></div>

#### function.Params

``` suneido
() => string
```

Returns a string containing the parameter list for the function.

Works for library functions, built-in functions, and blocks.

For example:

``` suneido
SendMessage.Params() => "(hwnd,msg,wParam,lParam)"
```

Useful for calltips in source code editors.

See also: [block.Params](<../Block/block.Params.md>)