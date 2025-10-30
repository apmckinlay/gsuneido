<div style="float:right"><span class="builtin">Builtin</span></div>

#### record.Observer

``` suneido
( observer )
```

Registers *observer *to be called whenever a member is modified.  
Multiple observers can be registered on the same record.
Observers are called in the same order they were registered.

Observers are also triggered by the invalidation of rules.

The observer call is equivalent to:

``` suneido
record.Eval(observer, member: member)
```

**Note**: this means that if the observer wants to receive the member its parameter <u>must</u> be called "member".

**Note**: If the observer is a method of a class, it will be called normally, and "this" will be the instance or class, <u>not</u> the record. (This is because when you get a reference to a method, the value you get is "bound" to the class or instance that you got it from.)

For example:

``` suneido
f = function (member) { Print(member $ " changed to " $ this[member]) }
r = new Record
r.Observer(f)
r.a = 1
r.b = 2
```

would print:

``` suneido
a changed to 1
b changed to 2
```

See also:
[record.RemoveObserver](<record.RemoveObserver.md>)