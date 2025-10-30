### Running Multiple Methods

Sometimes you want to execute multiple methods in a class. And you want to be able to easily add additional methods.

If you explicitly call each of the methods, then you are duplicating the method list, and to add a method you have to remember to add it in two places.

You could put all the code into one big function, but that is not as easy to read or to modify. Splitting the code into named functions also means you have a method name to help understand what the code is doing.

One approach to handling this is to loop through the members of the class. For example:

``` suneido
class
    {
    CallClass()
        {
        for m in .Members()
            if m.Has?('contrib_')
                this[m]()
        }
    contrib_foo()
        {
        ...
        }
    contrib_bar()
        {
        ...
        }
```

Of course, you can pass arguments or accumulate or check return values.

Normally, if you have a private method like "contrib_one" that is not called anywhere you will get a code warning. You can avoid the warning by making the methods public but it is better if they are private. The code checker specifically ignores methods named "contrib_..."

**Note:** As usual, the order that .Members() returns is undefined (hash map order). So if we need to run the methods in a specific order we can sort the members and name the members so they end up in the order we want.
<pre>
class
    {
    CallClass()
        {
        for m in .Members()<b>.Sort!()</b>
            if m.Has?('contrib_')
                this[m]()
        }
    contrib<b>1</b>_first()
        {
        ...
        }
    contrib<b>2</b>_second()
        {
        ...
        }
</pre>