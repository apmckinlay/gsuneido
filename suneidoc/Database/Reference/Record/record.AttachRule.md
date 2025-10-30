<div style="float:right"><span class="builtin">Builtin</span></div>

#### record.AttachRule

``` suneido
(member, callable)
```

Explicitly associate a rule with a field. The rule may be anything callable i.e. function, method, block, or class.

This is an alternative to automatically associating rules with fields by name (e.g. Rule_myfield with myfield).

**Note:** If the default value is removed from a record with `record.Set_default()`, then rules will no longer be automatically associated by name and <u>only</u> rules that are explicitly associated by record.AttachRule will be used.

For example:

``` suneido
x = Record()
Assert(x.foo is: "")
x.Set_default()
Assert({ x.foo } throws: "member not found")
x.AttachRule(#foo, function () { "bar" })
Assert(x.foo is: "bar")
```

See also: 
[Rules](<../../Rules.md>)