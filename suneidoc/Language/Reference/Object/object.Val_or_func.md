#### object.Val_or_func

``` suneido
(member) => value
```

If the member is a function, it will be called and the result returned. Otherwise the value of the member will be returned. For example:

``` suneido
#(a: 12, b: 34).Val_or_func("b") => 34

#(f: Date).Val_or_func("f") => #20231016.1951
```