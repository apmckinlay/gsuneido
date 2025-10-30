### Test

Abstract base class for *tests*.

Methods:
`CallClass()`
: run the test, give an alert with the results

`Run(quiet = false)`
: run the test, print the results

`RunDebug()`
: run the test, but don't catch exceptions so debugger will come up

`Setup()`
: empty default method

`MakeTable(schema) => name`
: create a table, for example:
``` suneido
name = .MakeTable("(name, phone) key(name)")
```

`MakeView(definition) => name`
: create a view, for example:
``` suneido
name = .MakeView("tables join columns")
```

`MakeLibrary(@records) => name`
: create a library, for example:
``` suneido
name = .MakeLibrary()
```

`MakeFile(@records) => name`
: Returns a name for a file but does not actually create it. The files will be deleted during teardown.

`AddTeardown(function)`
: register a teardown function (or method) which will then be called during Teardown

`SpyOn(target) => Spy`
: create a Spy on a target. target can be a global function, a class method or a global name string. The created spy will be cleaned up at the end of the current test function scope automatically no matter the test succeed or failed.

`WatchTable(tableName) => watchName`
: Creates a temporary trigger to save changes into the server memory, then you can use .GetWatchTable(watchName) to get new/modified records.
It also cleans up records created in the table when tearing down
For example:
``` suneido
log = .WatchTable("suneidolog")
Suneidolog('Error: testing')
logs = .GetWatchTable(log)
Assert(logs isSize: 1)
Assert(logs[0].sulog_message is: 'Error: testing')
```

`Teardown()`
: Runs teardown functions registered with AddTeardown.
: **Warning:** If you override Teardown, you must call super.Teardown()

See also:
[Using the Unit Testing Framework](<../../Cookbook/Using the Unit Testing Framework.md>)