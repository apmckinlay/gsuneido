### XmlRpcMap

A function that should return an object mapping names to functions, used by [XmlRpc](<XmlRpc.md>). For example:

``` suneido
function ()
    {
    return Object(
        "test.echo": function (@args) { return args }
        "examples.getStateName": function (i) { return States[i] }
        "system.memory": MemoryArena
        )
    }
```

The default XmlRpcMap in stdlib contains test.echo, test.throw, and examples.getStateName.

In your own library, to add to the stdlib XmlRpcMap (rather than replace it), call _XmlRpcMap, for example:

``` suneido
function ()
    {
    map = _XmlRpcMap()
    map["my.one"] = MyOne
    map["my.two"] = MyTwo
    return map
    }
```