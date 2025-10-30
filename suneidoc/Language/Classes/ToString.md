### ToString

If a instance has a `ToString()` method it will be used by [Display](<../Reference/Display.md>) (and therefore [Print](<../Reference/Print.md>) and other functions that use Display). It will also be used by string concatenation ($).

This can be very useful for debugging so you can see the contents of class objects.

Consider writing a ToString() method for your classes.

For example:

``` suneido
Point
class
    { 
    New(x, y)
        { 
        .x = x
        .y = y
        } 
    ...
    ToString() 
        { 'Point(x:' $ Display(.x) $ ', y: ' $ Display(.y) $ ')' }
    }

Point(12, 34)
    => Point(x: 12, y: 34)
```

Without the ToString method you would only see:

``` suneido
Point(12, 34)
    => Point()
```

Note: When possible, it is a good idea to make ToString's output match the syntax that would be used to create the instance, as in the Point example above.

**Note**: ToString only works on instances, not on the class itself.