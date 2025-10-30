<div style="float:right"><span class="builtin">Builtin</span></div>

### ServerEval

``` suneido
(function_name [, args ...]) => value
```

Calls the named function, passing it the supplied arguments, and returns the result. When running client-server, the function is called on the server.

For example:

``` suneido
ServerEval("QueryAll", "tables", limit: 100)
```

The function can be a class method, for example:

``` suneido
ServerEval("Dates.Begin")
    => #17000101
```

**Tip**: You can build the arguments in an object and use '@', for example:

``` suneido
ob = Object("QueryAll")
ob.Add(tablename)
ob.limit = limit
ServerEval(@ob)
```

**Note**: An instance method, e.g. "x.Func" will <u>not</u> work since "x" will not be available on the server.

Unlike [string.ServerEval](<String/string.ServerEval.md>), ServerEval() cannot be used to evaluate expressions like "Suneido.member". But it is better to put the code you want to run on the server into a function anyway.

**Warning**: It is not a good idea to ServerEval anything that could "block" (e.g. System, Sleep, RunPiped) since this will prevent the server from responding to other users' requests.