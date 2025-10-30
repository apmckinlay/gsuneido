<div style="float:right"><span style="font-weight:bold;font-variant:small-caps">Compiler Error</span></div>

### too many literals

A single function (or class method) is limited to 256 *literals* e.g. strings or object constants.

In most cases you can avoid this limit by placing large numbers of literals inside an object constant, which will then only count as a single literal.

Instead of:

``` suneido
a = "abra"
b = "cada"
c = "bra"
...
```

Use:

``` suneido
ob = #(
    a: "abra"
    b: "cada"
    c: "bra"
    ...)
```

**Note**: Integers from -32k to +32k (16 bit) are not counted as literals.