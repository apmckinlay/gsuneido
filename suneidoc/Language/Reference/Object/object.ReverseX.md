<div style="float:right"><span class="builtin">Builtin</span></div>

#### object.Reverse!

``` suneido
() => object
```

Reverse the order of the *list* members of the object (consecutive members starting at 0).  Has no effect on any named members.

**Note:** This modifies the object it is applied to, it does not create a new object.

For example:

``` suneido
Object(12, 34, 56, 78).Reverse!() => #(78, 56, 34, 12)
```

See also: [object.Sort!](<object.Sort!.md>)