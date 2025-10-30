<div style="float:right"><span class="builtin">Builtin</span></div>

### Construct

``` suneido
( class, suffix="" ) => object
( object, suffix="" ) => object
```

Creates an instance of a class.  class can be either an actual class or else a string that is combined with suffix to get the class name.

If an object is passed instead of just a class [0] is taken as the class, and the rest of the object (@+1) is passed as arguments.

For example:

``` suneido
Construct( Object ) => Object( )
Construct( "Object" ) => Object( )
Construct( #( Object, 1, 2, 3 ) ) => Object(1, 2, 3)
Construct( 'Button', 'Control' ) => ButtonControl( )
```
Construct is used primarily for user interface controls and report formats. The normal way to construct an instance of a class is to use "new" or call the class. (See 
[Classes Overview](<../Classes/Overview.md>))