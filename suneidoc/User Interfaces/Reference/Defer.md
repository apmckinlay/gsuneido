<div style="float:right"><span class="builtin">Windows Builtin</span></div>

### Defer

``` suneido
(block) => killer
```

Queues the block to run from the message loop, which means it will be deferred until after the current callback returns. The block is run from a shared 10ms (the minimum) Windows timer. Timers have low priority so if the UI is busy there may be delays.

killer.Kill() will remove the block from the pending queue **if it is still there**.


See also:
[Delay](<Delay.md>)
