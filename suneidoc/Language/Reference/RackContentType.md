### RackContentType

Middleware for [RackServer](<RackServer.md>).

Supplies default Content_Type for paths based on MimeTypes, for example:

|  |  | 
| :---- | :---- |
| js | application/javascript | 
| html | text/html | 
| css | text/css | 
| ico | image/x-icon | 


Middleware is specified with the RackServer with argument, e.g.

``` suneido
RackServer(app: MyApp, with: [RackContentType])
```

**Note:** To override RackContentType you must supply Content_Type, not Content-Type


See also:
[RackServer](<RackServer.md>),
[RackRouter](<RackRouter.md>),
[RackLog](<RackLog.md>),
[RackEcho](<RackEcho.md>),
[RackDebug](<RackDebug.md>)
