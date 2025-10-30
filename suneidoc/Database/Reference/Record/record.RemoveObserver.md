<div style="float:right"><span class="builtin">Builtin</span></div>

#### record.RemoveObserver

``` suneido
( observer )
```

Remove an observer from a record that was added with [record.Observer](<record.Observer.md>)

For example:

``` suneido
f = function (member) { Print(member $ " changed to " $ this[member]) }
r = new Record
r.Observer(f)
r.a = 1
r.RemoveObserver(f)
r.b = 2
```

would print:

``` suneido
a changed to 1
```

Notice that it does <u>not</u> print "b changed to 2".