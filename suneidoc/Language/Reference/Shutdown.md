### Shutdown

``` suneido
(alsoServer = false)
```

Exits from Suneido. If **alsoServer** is true, it first uses ServerEval to call Shutdown *on the server*. When Shutdown is run on a server, it delays its exit by 5 seconds to allow the client time to Exit properly.

See also:
[Exit](<Exit.md>)