### Setup Before Base Constructors

Normally, base constructors (New's) are done first.  Normally, this is the right behaviour. However, occasionally, you need to do some set up before the base constructor is called. You can do this by using a function call as an argument to an explicit call to the base constructor.

For example:

``` suneido
class
    {
    New(arg)
        {
        super(.setup(arg))
        ...
        }
    setup(arg)
        {
        // your pre-base constructor set up goes here
        return arg
        }
```