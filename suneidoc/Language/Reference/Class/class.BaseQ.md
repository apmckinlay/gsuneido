<div style="float:right"><span class="builtin">Builtin</span></div>

#### class.Base?

``` suneido
(class) => true or false
```

Returns true if the class or instance is derived from the specified class (possibly indirectly), false if not.

For example:

``` suneido
HorzControl.Base?(Group) 
    => true

Stack().Base?(Stack) 
    => true
```

A class will return true for itself:

``` suneido
Stack.Base?(Stack)
    => true
```

**Note:** The following sequence of events can be confusing:

0.	`x = new MyClass`
1.	`x.Base?(MyClass) => true`
2.	`Unload("MyClass")`
3.	<code>x.Base?(MyClass) => <u>false</u></code>


You normally are not Unload'ing yourself, but editing the record in Library View or using or unusing libraries will have the same effect.

The reason for this is that the instance of the class (stored in x above) remains linked to the old version of MyClass, but you are doing the test with the new version.

See also: [class.Base](<class.Base.md>)