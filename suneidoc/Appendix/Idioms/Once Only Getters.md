### Once Only Getters

``` suneido
getter_name()
    {
    ... 
    return .name = ... 
    }
```

The first time `.name` is referenced, its getter method will be evaluated, setting a member of the same name. Subsequent references to `.name` will use the member, and will not execute the method. 'name' may be either private or public (capitalized e.g. `Getter_Name`).

**Note:** If the once method is returning an object, consider making it readonly. Since the result will be shared, you normally don't want it modified. For example:

``` suneido
getter_columns()
    {
    ob = Object('name')
    for (i = 0; i &lt 10; ++i)
        ob[i] = "col" $ i
    ob.Set_readonly()
    return .columns = ob
    }
```