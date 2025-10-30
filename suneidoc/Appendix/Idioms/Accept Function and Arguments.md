### Accept Function and Arguments

``` suneido
function (@args)
    {
    fn = args[0]
    ...
    fn(@+1args)
    ...
    }
```

Allow zero or more arguments for a function to be passed along with the function.

Often you want to pass an argument to the function in addition to any passed in:

``` suneido
function (@args)
    {
    fn = args[0]
    ...
    args[0] = x
    fn(@args)
    ...
    }
```