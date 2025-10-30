<div style="float:right"><span class="builtin">Builtin</span></div>

### Unload

``` suneido
(string = false)
```

Removes the definition of the named global. Used by Library View to remove old values when definitions are changed.

If called without a name i.e. Unload() it unloads all library globals from memory, similar to the effect of [Use](<Use.md>) or [Unuse](<Unuse.md>). **Note**: This should not be done lightly (or frequently) since it can have significant impact on performance.