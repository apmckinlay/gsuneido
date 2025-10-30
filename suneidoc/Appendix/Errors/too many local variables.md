<div style="float:right"><span style="font-weight:bold;font-variant:small-caps">Compiler Error</span></div>

### too many local variables

The normal solution to this is to split up your function, often a good idea for readability anyway. Another alternative is to use an object to store a collection of variables.

Instead of:

``` suneido
a = "abra"
b = "cada"
c = "bra"
...
```

Use:

``` suneido
ob = Object()
ob.a = "abra"
ob.b = "cada"
ob.c = "bra"
...
```