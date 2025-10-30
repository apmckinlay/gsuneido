### Same?

``` suneido
(x, y) => true/false
```

Returns true if the two values refer to the same object. In Suneido, is or == compare values so two objects with the same contents will be equal, even though they are different objects.

For example:

``` suneido
x = Object(123, m: 456)
y = Object(123, m: 456)
x is y
    => true
Same?(x, y)
    => false
y = x
Same?(x, y)
    => true
```

**Note**: Same? is for special uses, normal code should just use is/==

See also: [Hash](<Hash.md>)