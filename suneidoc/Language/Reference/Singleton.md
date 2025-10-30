### Singleton

An abstract class used to derive "singleton" classes (classes that support a single instance.

For example:

**`Logger`**
``` suneido
Singleton
    {
    New()
        { .log = Object() }
    Append(msg)
        { .log.Add(msg) }
    Log()
        { return .log }
    }
```

Singleton's are accessed via CallClass, for example:

``` suneido
Logger().Append("hello world")
```

Singleton provides a Reset method that destroys the instance so it will be recreated the next time it is used.

``` suneido
Logger().Reset()
```

**Note:** Singleton does not stop you from creating multiple instances.

See also: [Implementing the Singleton Pattern](<../../Cookbook/Implementing the Singleton Pattern.md>) and [CallClass for Singleton](<../../Appendix/Idioms/CallClass for Singleton.md>)