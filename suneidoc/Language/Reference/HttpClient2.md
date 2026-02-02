<div style="float:right"><span class="builtin">Builtin</span></div>

### HttpClient2

``` suneido
(method, url, content = '', header = #(), timeout = 60, block = false) 
 => object(content, header) or nil
```

Makes an HTTP request using the specified method (e.g., "GET", "POST", "PUT", "DELETE").

HttpClient is a thin wrapper around [Go net/http](<https://pkg.go.dev/net/http>)

**Parameters:**

`method`
: The HTTP method as a string (e.g., "GET", "POST")

`url`
: The URL to request as a string. Both http and https are handled.

`content`
: Optional request body. Can be a string or a function. If a function, it will be called repeatedly with a buffer size and should return chunks of data, or an empty string to signal EOF. To convert to push rather than pull, use [Pipe](<Pipe.md>)

`header`
: Optional object containing HTTP headers. Keys are header names (underscores are converted to hyphens). Special handling for Content-Length.

`timeout`
: Request timeout in seconds (default 60)

`block`
: Optional function to process the response.

**Return Value:**

If no block is provided:
- Returns an object with `content` (response body) and `header` (response headers) fields

If a block is provided:
- Calls the block with a response object and returns nil

#### Response Object Methods (when using block):

`Header() => string`
: Returns the response headers as a formatted string including the status line

`Status() => number`
: Returns the HTTP status code (e.g., 200, 404)

`Read(n) => string | false`
: Reads up to n bytes from the response body. Returns a string of the bytes read, false on EOF, or empty string if no data available. Maximum read size is 64kb.

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

**Example:**

```suneido
// Simple GET request
result = HttpClient2("GET", "https://api.example.com/data")
Print(result.content)

// POST with headers
result = HttpClient2("POST", "https://api.example.com/submit",
    "request data",
    #(Content_Type: "application/json"))

// Using a block to process response
HttpClient2("GET", "https://example.com/largefile",
    block: function (resp)
        {
        Print("Status:", resp.Status())
		File("output.dat", "w")
			{|f|
			resp.CopyTo(f)
			}
        })
```