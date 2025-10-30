### Http

``` suneido
(method, addr, content = false, fromFile = '', toFile = '', 
	header = #(), timeout = 60, timeoutConnect = 60, asyncCompletion = false)
	=> [header: header, content: content]

.Get(@args) => content

.Post(@args) => content

.Put(@args) -=> content
```

Perform an HTTP request.

**Note**: Http.Get/Post/Put check the result code and throw an exception if it is not 200 OK. They are preferred over `Http(method` since it is a good idea to always check the result code.

timeout and timeoutConnect are in seconds.

If asyncCompletion is given the request will be done in a separate thread and asyncCompletion will be called in that thread. If it needs to do UI it can use [Defer](<../../User Interfaces/Reference/Defer.md>).

For example:

``` suneido
Http.Get("www.suneido.com")
```

See [SplitAddr](<SplitAddr.md>) for a description of the address format.