### Build Arguments for Base

``` suneido
Base
    {
    New(@args)
        {
        super(@.makeargs(args))
        }
    makeargs(args)
        {
        ...
        return args_for_base
        }
    ...
```

Since a call to a base class constructor must be the first statement in the New function, if you need to build its arguments you have to do so in a separate function.