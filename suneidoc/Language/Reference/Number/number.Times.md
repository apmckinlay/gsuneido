#### number.Times

``` suneido
(callable)
```

Does:

``` suneido
for (i = 0; i < number; ++i)
    callable()
```

callable can be anything that can be called e.g. block, function, class, or instance.

For example:

``` suneido
3.Times()
    { Print("hello") }
    =>  "hello"
        "hello"
        "hello"
```