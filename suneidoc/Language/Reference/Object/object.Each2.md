#### object.Each2

``` suneido
(callable) => this
```

Does:

``` suneido
for member in object.Members()
    callable(member, object[member])
```

callable can be anything that can be called e.g. block, function, class, or instance.

For example:

``` suneido
#(1, 2, a: 3, b: 4).Each2({|m,v| Print(m, '=', v) })
    =>  0 = 1
        1 = 2
        b = 4
        a = 3
```

**Note:** The order of named members is unpredictable.

See also: [object.Each](<object.Each.md>)