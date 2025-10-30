### Default

When a member function is called, if the member does not exist, then a member function called "Default" is called with the member name as the first argument and the remaining arguments from the original call.

For example, to forward unrecognized method calls to another object:

``` suneido
class
    {
    ...
    Default(@args)
        {
        return .other[args[0]](@+1args);
        }
```