<div style="float:right"><span class="builtin">Builtin</span></div>

#### object.Readonly?

``` suneido
() => true or false
```

Returns true if the object is readonly, false if not.
An object may be readonly because it originated from a literal (i.e. `#(...)`)
or because Set_readonly was used on it.

For example:

``` suneido
#().Readonly?() => true
Object().Readonly?() => false
Object().Set_readonly().Readonly?() => true
```