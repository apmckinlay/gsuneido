<div style="float:right"><span class="builtin">Builtin</span></div>

#### socketClient.Writeline

``` suneido
(string)
```

Write the string, followed by "\r\n", to the socket connection.
(i.e. The string would normally *not* end in "\r\n".)

An exception will be thrown if the timeout is exceeded.