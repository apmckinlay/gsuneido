<div style="float:right"><span class="builtin">Builtin</span></div>

### DoWithoutTriggers

``` suneido
(tables, block)
```

Executes the block with the triggers for the specified tables disabled.

This is useful for mass updates when you know the triggers aren't necessary.

For example:

``` suneido
DoWithoutTriggers(#(mytable))
    {
    QueryApply('mytable', update:)
        { |x|
        x.name = Capitalize(x.name)
        x.Update()
        }
    }
```

**Note:** DoWithoutTriggers only works standalone or on the server. On a client it has no effect. If necessary you can use string.ServerEval to run a function that uses DoWithoutTriggers (so that it runs on the server where it will work), but be aware that this will disable the triggers for all clients.