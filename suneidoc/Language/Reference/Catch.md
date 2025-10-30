### Catch

``` suneido
( callable ) => return value or exception
```

Calls the supplied block, function, object, or class within a try-catch. If an exception occurs, the exception is returned.

For example:

``` suneido
func = function () { return 123 + 456 }
Catch(func)
    => 579

func = function () { return j }
Catch(func)
    => "uninitialized variable: j"

Catch(){ Database("drop mytable") }
    => "nonexistent table: mytable"
```

The last example is written with a block immediately following the Catch call so it will be interpreted as an argument to the call.