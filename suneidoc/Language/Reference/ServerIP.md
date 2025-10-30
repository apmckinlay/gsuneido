<div style="float:right"><span class="builtin">Builtin</span></div>

### ServerIP

``` suneido
() => string
```

Returns the IP address used if running client-server, otherwise "".

For example:

``` suneido
ServerIP()
    => "127.0.0.1"
```

See also:
[Server?](<Server?.md>),
[ServerPort](<ServerPort.md>)