## Implementing the Singleton Pattern

**Category:** Coding

**Ingredients**

-	the special 
	[CallClass](<../Language/Classes/CallClass.md>) method
-	the global Suneido object
-	the 
	[object.Base()](<../Language/Reference/Object/object.Base.md>) method


**Problem**

You want a single instance of a class that is available globally.

**Recipe**

Unlike other object oriented languages like Java or C++, Suneido does not allow static data members of classes. A class is a constant, it can't contain data members (other than other constants). The only place Suneido provides for global variables is the "Suneido" object. So to implement the Singleton pattern we store the instance in a member of the Suneido object. Here is the code:

*`Singleton`*
``` suneido
class
    {
    CallClass()
        {
        if not Suneido.Member?(.name())
            Suneido[.name()] = new this
        return Suneido[.name()]
        }
    Reset()
        {
        Suneido.Delete(.name())
        }
    name()
        {
        c = .Base() is Singleton ? this : .Base()
        return Display(c).BeforeFirst(' ')
        }
    }
```

Rather than use a static method (such as "Instance()" to access the instance we can use the CallClass method. This method is used if the class is "called". This allows us to access the singleton as simply:

``` suneido
MySingleton().MyMethod()
```

The CallClass method creates the instance if it doesn't exist yet.

The Reset method allows you to destroy the instance (so it will be recreated).

The name of the class is used as the member name in the Suneido object. If name() is called on the instance, it uses the Base() method to get the class.

Singleton is an abstract base class, to use it you must define your own class that inherits from Singleton. For example:

*`Logger`*
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

We could then use it like:

``` suneido
    Logger().Append("starting")
    ...
    Logger().Append("ending")
    ...
    Print(Logger().Log())
    Logger().Reset()
```
**Unit Test**
Here is a unit test for Singleton:

*`Singleton_Test`*
``` suneido
Test
    {
    Test_main()
        {
        c = Singleton_TestClass
        Assert(c() is: Suneido.Singleton_TestClass)
        Assert(c().N is: 0)
        c().Inc()
        Assert(c().N is: 1)
        c.Reset()
        Assert(not Suneido.Member?('Singleton_TestClass'))
        Assert(c().N is: 0)
        }
    Teardown()
        {
        Suneido.Delete('Singleton_TestClass')
        }
    }
```

Because Singleton depends on the class name, we have to use a separately defined test singleton:

*`Singleton_TestClass`*
``` suneido
Singleton
    {
    New()
        { .N = 0 }
    Inc()
        { ++.N }
    }
```

**Uses**

[Singleton](<../Language/Reference/Singleton.md>) is now used for the plugin manager ([Plugins](<../Language/Reference/Plugins.md>)) and for the new [AccessControl](<../User Interfaces/Reference/AccessControl.md>) lock manager.

**See Also**

Design Patterns by Gamma, Helm, Johnson, and Vlissides includes the Singleton pattern.

Head First Design Patterns by Freeman & Freeman covers the pattern in a more "friendly" style.

Appendix > Idioms > [CallClass for Singleton](<../Appendix/Idioms/CallClass for Singleton.md>)