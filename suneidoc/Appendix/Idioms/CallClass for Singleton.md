### CallClass for Singleton

**`MySingleton`**
``` suneido
class
    {
    CallClass()
        {
        if not Suneido.Member?("MySingleton")
            Suneido.MySingleton = new MySingleton
        return Suneido.MySingleton
        }
    ...
    }
```

Then, to access the singleton you use MySingleton() e.g. MySingleton().Func()

A normal method could be used instead of CallClass e.g. Instance, but this makes using it more awkward e.g. MySingleton.Instance().Func()

It is conventional to use the same name for the class as the Suneido member. Often a short name makes using it nicer. But, as with all global names, try to pick a name that is not likely to clash with other people's names.

Used by [Plugins](<../../Language/Reference/Plugins.md>).