<div style="float:right"><span class="builtin">Builtin</span></div>

#### string.MapN

``` suneido
(n, block) => string
```

Calls block(substr) for each substring of length n and returns the results concatenated together.

For example:

``` suneido
"hello world".MapN(2, { it.Capitalize() })
    => "HeLlO WoRlD"
```

The last substring may be less than n characters.

If the string is "" then an empty object is returned.

string.MapN may also be used for side effects, for example:

``` suneido
"hello world".MapN(2, { Print(it); "" })
    =>  he
        ll
        o 
        wo
        rl
        d
```

When using MapN for side effects only, it is best to make the value of the block "" (as in the above example).

See also:
[string.Map](<string.Map.md>),
[string.Divide](<string.Divide.md>)