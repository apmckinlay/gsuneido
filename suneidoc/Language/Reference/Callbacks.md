<div style="float:right"><span class="builtin">Builtin</span></div>

### Callbacks

``` suneido
( ) => object
```

Returns a list of the callbacks currently active.

A good way to use this is to Inspect the result:

``` suneido
Inspect(Callbacks())
```

This is primarily useful as a debugging tool to help ensure that callbacks are being cleared.
WndProc.NCDESTROY clears window procedure callbacks,
but other callbacks, such as timers, must be cleared manually.
If callbacks are not cleared, they will tie up memory and resources unecessarily.