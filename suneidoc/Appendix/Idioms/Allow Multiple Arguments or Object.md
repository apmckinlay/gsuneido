### Allow Multiple Arguments or Object

``` suneido
function (@args)
    {
    if args.Size() is 1 and args.Member?(0) and Object?(args[0])
        args = args[0]
```

Allow a function to be called with either a list of arguments:

``` suneido
fn( 1, 2, 3 )
```

or an object:

``` suneido
ob = #( 1, 2, 3 )
fn( ob )
```