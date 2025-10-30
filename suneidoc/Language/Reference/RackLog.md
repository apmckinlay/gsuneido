### RackLog

A [RackServer](<RackServer.md>) middleware component that logs requests.

It uses [Rlog](<Rlog.md>) to log the method, path, and response code plus the environment (excluding body, cookie, and auth_token) to rack#.log

For example:

``` suneido
RackServer(app: RackEcho, with: [RackLog])
```

The output looks like:

``` suneido
20160103.121547652 GET /wiki => 200 OK (1ms) #(host: "127.0.0.1", dnt: "1", accept: "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8", cache_control: "max-age=0", remote_user: "127.0.0.1", method: "GET", version: 1.1, queryvalues: #(), upgrade_insecure_requests: "1", path: "/wiki", request: "GET /wiki HTTP/1.1", query: "", connection: "keep-alive", user_agent: "Mozilla/5.0 (Windows NT 10.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/47.0.2526.106 Safari/537.36", accept_encoding: "gzip, deflate, sdch", accept_language: "en-US,en;q=0.8")
```


See also:
[RackServer](<RackServer.md>),
[RackRouter](<RackRouter.md>),
[RackContentType](<RackContentType.md>),
[RackEcho](<RackEcho.md>),
[RackDebug](<RackDebug.md>)
