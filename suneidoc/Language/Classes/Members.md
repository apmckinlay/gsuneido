### Members

Within classes, member names starting with a lower case letter are *private*. Member names starting with an upper case letter are *public*.

Private members are implemented by internally prefixing them with their class name.
<pre>
MyClass
class
    {
    <b>foo:</b> 123
    Bar()
        {
        return <b>.foo</b>
        }
    }
</pre>

Internally, both of the foo's will become "MyClass_foo".

``` suneido
MyClass.Members() => #("MyClass_foo", "Bar")
```

So access from outside the class (e.g. MyClass.foo) will throw an exception (member not found).

Private members can still be accessed from outside the class by using the internal name (e.g. x.MyClass_foo), but this is <u>not</u> recommended except for tests.