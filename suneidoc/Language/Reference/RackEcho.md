### RackEcho

A simple [RackServer](<RackServer.md>) app that returns its environment with one member per line.

For example:

``` suneido
RackServer(app: RackEcho)
```

Produces a web page like:

``` suneido
host: "127.0.0.1"
dnt: "1"
accept: "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8"
remote_user: "127.0.0.1"
method: "GET"
version: 1.1
upgrade_insecure_requests: "1"
queryvalues: #()
body: ""
query: ""
request: "GET /echo HTTP/1.1"
path: "/echo"
accept_encoding: "gzip, deflate, sdch"
user_agent: "Mozilla/5.0 (Windows NT 10.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/47.0.2526.106 Safari/537.36"
connection: "keep-alive"
accept_language: "en-US,en;q=0.8"
```


See also:
[RackServer](<RackServer.md>),
[RackRouter](<RackRouter.md>),
[RackContentType](<RackContentType.md>),
[RackLog](<RackLog.md>),
[RackDebug](<RackDebug.md>)
