#### object.Each

``` suneido
(callable) => this
```

Does:

``` suneido
for value in object
    callable(value)
```

callable can be anything that can be called e.g. block, function, class, or instance.

For example:

``` suneido
BookTables().Each(Print)
    =>  imagebook
        mybook
        suneidoc
```

**Note:** The order of named members is unpredictable.

See also: [object.Each2](<object.Each2.md>)