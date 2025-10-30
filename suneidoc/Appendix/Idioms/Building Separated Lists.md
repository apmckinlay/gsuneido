### Building Separated Lists

A common task is to create a string listing several items with separators between them. It's often easier to add the separator after every item, and then at the end, remove the trailing one.

``` suneido
s = "";
for (x in list)
    s $= Fn(x) $ ","
s = s[..-1]
```

An alternative is to build an object and then use Join.

``` suneido
ob = Object()
for (x in list)
    ob.Add(Fn(x))
s = ob.Join(',')
```

Or, for this example, it's even simpler with [object.Map](<../../Language/Reference/Object/object.Map.md>):

``` suneido
s = list.Map(Fn).Join(',')
```