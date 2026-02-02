### RackServer

``` suneido
(name = "Rack Server", port = 80, exit = false, app, with = #())
```

An HTTP [SocketServer](<SocketServer.md>) based on [Ruby Rack](<http://rack.rubyforge.org/doc/SPEC.html>) and [Python WSGI](<http://www.python.org/dev/peps/pep-3333>) designs. RackServer supports HTTP 1.1 persistent connections whereas HttpServer does not.

If an exception is thrown it will be logged with SuneidoLog and a '500 Internal Server Error' will be returned. Socket related exceptions are not logged.

RackServer currently does not support streaming input or output. Because it reads the entire request body into memory it will return '413 Request Entity Too Large' if the request Content-Length is too large.

#### Application

The supplied app is called with (env: env) which contains:
`remote_user`
: From 
[socketServer.RemoteUser](<SocketServer/socketServer.RemoteUser.md>)()

`request`
: The complete request line

`method`
: The method from the request in upper case e.g. "GET"

`path`
: The portion of the URI before a question mark ('?')

`query`
: The portion of the URI after a question mark ('?')

`queryvalues`
: The query string converted to an object using 
[UrlDecodeValues](<UrlDecodeValues.md>)

`version`
: The HTTP version from the request line as a number, e.g. 1.1

`body`
: The content of the request

plus the headers from the request as interpreted by InetMesg.HeaderValues.

The app should return an object with three members - the response code, an object containing response headers, and a content string. The response header members will be translated from underscores to dashes so that e.g. you can write `#(Content_Type: 'text/html')` instead of `#('Content-Type': 'text/html')`

As a shortcut, returning just a string will default to a response code of "200 OK" and no extra headers.

Returning an object with two members will be taken as the response code and the content, with no extra headers.

The simplest app is a function returning a string, for example:

``` suneido
function () { return "hello world" }
```

which is equivalent to:

``` suneido
function (env/*unused*/) { return #("200 OK", (), "hello world") }
```

Note: Since this simple app does not look at the request, it will return the same result for any method or url.

Run it with:

``` suneido
RackServer(app: function () { return "hello world" })
```

Note: the app argument must be named since we are not passing the name, port and exit arguments.

#### Middleware

Middleware can be specified using the "with" parameter. For example:

``` suneido
RackServer(app: MyApp, with: [RackContentType, RackLog])
```

The middleware will be applied in the order listed, in this example RackLog would be executed last (outermost).

The simplest "do nothing" middleware is:

``` suneido
class
    {
    New(.app)
        { }
    Call(env)
        {
        return (.app)(env: env)
        }
    }
```

A middleware component can do things before or after the app that it wraps. It can alter the request environment before passing it to the app, or it can alter the response returned by the app. It also has the option of handling the request itself and not calling the app that it wraps at all.

Middleware that wants to standardize a result can call RackServer.ResultOb(result) which handles the result shortcuts and returns a standard three member object.


See also:
[RackRouter](<RackRouter.md>),
[RackContentType](<RackContentType.md>),
[RackLog](<RackLog.md>),
[RackEcho](<RackEcho.md>),
[RackDebug](<RackDebug.md>)
