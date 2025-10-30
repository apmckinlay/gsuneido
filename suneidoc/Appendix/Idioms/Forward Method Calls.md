### Forward Method Calls

To forward unrecognized method calls to another object:

``` suneido
class
    {
    ...
    Default(@args)
        {
        return .other[args[0]](@+1 args)
        }
```