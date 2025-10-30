### XmlRpc

A simple implementation of XML-RPC.

CallClass is the "server" portion. HTTP messages to HttpServer with a URL of XmlRpc will be processed here. Exceptions are returned as XML-RPC "faults". XML-RPC values are converted to Suneido values as follows

|  |  | 
| :---- | :---- |
| **XML-RPC** | **Suneido** | 
| boolean | boolean | 
| int | number | 
| double | number | 
| string | string | 
| struct | object | 
| array | object | 
| dates | not implemented yet | 
| base64 | not implemented yet | 


Call is the "client" portion. Use it to make an XmlRpc call, for example:

``` suneido
XmlRpc.Call('xmlrpc.usefulinc.com/demo/server.php', 'examples.getStateName', 41)
```

The first argument is the URL, the second is the method name. Any additional arguments are passed to the method. Fault result cause an exception to be thrown. Suneido values are converted as follows:

|  |  | 
| :---- | :---- |
| **Suneido** | **XML-RPC** | 
| boolean | boolean | 
| numbers that are integers | int | 
| other numbers | double | 
| string | string | 
| objects with named members | struct | 
| objects with only un-named members | array | 
| dates | not implemented yet | 


Calls are implemented using [HttpPost](<HttpPost.md>).

Encodes and decodes &amp; &lt; &gt; &quot; character entities.