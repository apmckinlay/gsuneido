### RackDebug

A [RackServer](<RackServer.md>) middleware component that logs requests.

``` suneido
RackServer(app: function () { throw "error" }, with: [RackDebug])
```

Would both Print and return something like: (with a response code of "500 Internal Server Error"

``` suneido
EXCEPTION: error
 /* function */
  RackDebug.Call /* stdlib method */
    RackServer.Run /* stdlib method */
```


See also:
[RackServer](<RackServer.md>),
[RackRouter](<RackRouter.md>),
[RackContentType](<RackContentType.md>),
[RackLog](<RackLog.md>),
[RackEcho](<RackEcho.md>)
