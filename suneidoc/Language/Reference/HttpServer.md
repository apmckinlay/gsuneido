<div style="float:right"><span class="builtin">Builtin</span></div>

### HttpServer

``` suneido
(port, app)
```

A simple HTTP server that listens on the given TCP port and forwards each
request to a Suneido handler function (the `app`). It is intended to replace
[RackServer](<RackServer.md>). Each connection is handled in its own thread via Go's
`net/http` package.

HttpServer is a thin wrapper around [Go net/http](<https://pkg.go.dev/net/http>)

If the listener encounters an error it panics with the underlying Go
message. Exceptions raised while handling a request return **500 Internal Server
Error** and are logged with a Go stack trace.

The `app` receives a single named argument `env`. It exposes commonly used
request metadata plus helpers for working with the body and response.

#### Env members

`method`
:	HTTP verb in uppercase, e.g. `"GET"`

`path`
:	Path portion of the URL (before `?`)

`query`
:	Raw query string (after `?`)

`queryvalues`
:	Object decoded from the query string. Repeated keys become arrays; keys with
	no value appear as members without a value.

`body`
:	Entire request body as a string. The body is read lazily. 
	Request bodies are read into memory when accessed through `env.body`.
	Use `Read` or `CopyTo` for large uploads to avoid holding the entire content in memory.

Any HTTP header can be retrieved by its name e.g. `env.Content_Length`

#### Env Methods

`Read(n) => string or false`
:	Reads up to n bytes from the request body content.
	Prefer env.Read to env.Body

`Write(string)`
:	Writes the string to the response body content.
	Either Content-Length is added or chunked mode is used, automatically.
	For short response body content, just return it from the app, Write is unnecessary.
	If Write is used then the return value from the `app` is ignored.

`WriteHeader(status, header = #())`
:	Writes the response header. 
	This is optional, it is only necessary if using `Write` **and** 
	you need a different status than 200 or to add headers.
	`status` is a numeric HTTP status code (e.g. `200`)
	NOTE: WriteHeader must be called before Write

`CopyTo(dest, nbytes = false) => ncopied`
:	Up to nbytes (or everything) is read from this and written to a suitable destination. 
	e.g. [File](<File.md>),
	[Pipe](<Pipe.md>),
	[HttpClient2](<HttpClient2.md>),
	[HttpServer](<HttpServer.md>),
	[RunPiped](<RunPiped.md>),
	[SocketClient](<SocketClient.md>),
	[SocketServer](<SocketServer.md>).
	Prefer `CopyTo` when applicable since it is more efficient.

#### Responses

The `app` can respond to a request in two ways:

-	Return one of the following values:
	-	A **string** - treated as the body with status `200` and no custom headers.
	-	A **two element object** - `#(status, body)` - no custom headers.
	-	A **three element object** - `#(status, headers, body)`
		where `headers` is an object of header names to values.

-	Call `env.WriteHeader` / `env.Write` - useful for streaming large content.
	In which case the return value is ignored.

`status` is a numeric HTTP status code (e.g. `200`)

#### Headers

When referencing headers, underscores are translated to dashes.
This is to avoid needing to quote them.

e.g. env.Some_Header instead of env["Some-Header"]
or #(Some_Header: "value") instead of #("Some-Header": "value")

Single values are strings, multiple values (duplicate headers) are list objects.

#### Example

``` suneido
MyHttpApp = function (env) {
    if env.path is: '/hello'
        return #(200, #(Content_Type: 'text/plain'), 'hello world')
    return #(404, 'Not Found')
}

Thread({ HttpServer(8080, app: MyHttpApp) })
```