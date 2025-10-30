<div style="float:right"><span class="builtin">Windows Builtin</span></div>

### Delay

``` suneido
(delayMs, block) => killer
```

Runs the block from the main thread message loop after the specified delay.

**Note**: Delay can only be used from the GUI thread.

killer.Kill() will prevent the block from being run.


See also:
[Defer](<Defer.md>)
