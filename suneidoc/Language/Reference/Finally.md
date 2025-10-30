<div style="float:right"><span class="builtin">Builtin</span></div>

### Finally

``` suneido
(main_block, final_block) => result
```

Like try-finally in e.g. Java. Runs the main block and then the final block. The final block is run even if the main block throws an exception.

If the main block does not throw and the final block does, then the final block exception will propagate. If the main block throws an exception, and then the final block also throws an exception, the final block exception will be ignored and the main block exception will propagate.

Returns the result of the main block (i.e. the value of the last statement) if there are no uncaught exceptions.

See also: [Exceptions](<../Exceptions.md>)