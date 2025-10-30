### Getter Methods

If a non-existent member is accessed, Suneido will look first for a method named "Getter_", and then for a method named "Getter_" $ name (for public members) or "getter_" $ name (for private members). If a *getter* method is found it will be called and the value of the member will be whatever the method returns. Getter_ is called with the name of the member being accessed. This allows members to be implemented as data or functions without altering client code.

For example:

``` suneido
myClass = class { i: 0, Getter_Num() { return .i++; } }
x = new myClass
Assert(x.Num is 0)
Assert(x.Num is 1)
```

Getters can be used to defer calculation of a value until it is first used, and to cache that value for further use:

``` suneido
myClass = class { N: 0, Getter_Num( ) { ++.N; return .Num = 123; } }
x = new myClass
Assert(x.Num is 123)
Assert(x.Num is 123)
Assert(x.N is 1) // Getter_Num() should only be called once
```

Because we set the member, the getter will no longer be called. (Because actual members take precedence over getters.)

**Note**: To work correctly with [object.GetDefault](<../Reference/Object/object.GetDefault.md>), Getter_ should return *nothing* for "missing" members. For example:

``` suneido
Getter_(member)
    {
    if .data.Member?(member)
        return .data[member]
    else
        return
    }
```

Notice that if the member is not available, it is not returning false or "" or any other value, it is using a return with <u>no</u> value.

**Note:** For private members, the member name passed to Getter_ will be "privatized", i.e. it will be prefixed with the class name.

**Note:** Previously (prior to 2019), getters were named "Get_" and "get_" but these were not specific/unique enough.

See also:
[Once Methods](<../../Appendix/Idioms/Once Methods.md>)