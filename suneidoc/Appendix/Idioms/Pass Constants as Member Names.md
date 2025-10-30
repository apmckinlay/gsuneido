### Pass Constants as Member Names

Constants are commonly defined in groups in objects. Most of the Windows constants are defined this way. For example, button control states are defined in an object called BST:

``` suneido
#(
UNCHECKED:      0x0000
CHECKED:        0x0001
INDETERMINATE:  0x0002
PUSHED:         0x0004
FOCUS:          0x0008
)
```

At first glance it seems simplest just to pass the constant values e.g. BST.CHECKED to functions or methods. However, it's often preferable to pass the member name instead. Because something like BST.CHECKED is an expression, not a constant, it cannot be used as a default parameter value.

``` suneido
function (state = BST.CHECKED) // ILLEGAL!
```

But if you pass member names instead, then you can write:

``` suneido
function (state = "CHECKED") // OKAY
```

Another advantage is that if you pass an invalid member name the error will be detected when the function/method does:

``` suneido
BST[state]
```

If you just pass numeric constants, there is no equivalent detection of bad values.