## Containers

The two basic methods of laying out multiple formats are Horz and Vert.  Horz arranges the items one beside the other horizontally, whereas Vert arranges the items one below the other vertically.

The items supplied as arguments to Params are automatically placed into a Vert, so:

``` suneido
Params(#(Text, "Hello world!"), #(Text, "Goodbye."))
```

will print "Goodbye." immediately below "Hello world!".  In addition to the normal Vert behavior, Report will automatically start a new page if the next item will not fit on the page.  You can also force a new page with a "pg" item.  For example,

``` suneido
Params(#(Text "Hello"), 'pg', #(Text "Goodbye"))
```

To place items side by side, use an Horz.  For example:

``` suneido
Params(#(Horz, (Text, "Hello"), (Text, "goodbye")))
```

Horz aligns text vertically on it's baseline, rather than just the top or bottom of the text, for example:

align <span style="font-size: xxx-large;">align</span>

Horz's and Vert's can be nested, for example:

``` suneido
#( Horz ( Vert (Text one) (Text two) ) ( Vert (Text three) (Text four) ) )
```

to produce:

``` suneido
onethree
twofour
```

Horz's and Vert's also allow you to override their default positioning of each item immediately following the last.  For example:

``` suneido
Params( #( Text "0,0" ),
    #( Horz ( Text "1,1" x: 1 ) y: 1 ),
    #( Horz ( Text "2,2" x: 2 ) y: 2 ) )
```

would produce:

``` suneido
0,0



          1,1



                    2,2
```

You can only set the x position for items within Horz's and the y position for items within Vert's.  If the specified position is smaller than the default position, it is ignored.  Positions are in inches and are relative to their container (not to the page).

Horz and Vert can also distribute data to their items.  If a data value is supplied it will be passed on to each of the format items.  If a format item has a "field" member, then the value of that field from the data object will be passed to the format rather than the entire data object.  For example:

``` suneido
#( Horz ( Text field: name ) data: ( name: "joe" ) )
```

is equivalent to:

``` suneido
#( Horz ( Text "joe" ) )
```

Using a field name instead of a format will look up that field's format, 
so you can write, for example:

``` suneido
#( Horz name age, data: ( name: "Joe", age: 23 ) )
```

A quick and easy way to output several fields is to use the Hfields and Vfields formats 
which output fields along with their prompts from the data dictionary.  For example:

``` suneido
#( Vfields name, age, data: ( name: "Joe", age: 23 ) )
```

would produce

``` suneido
Name: Joe
 Age: 23
```

These features are most useful in conjunction with
[QueryFormat](<Reference/QueryFormat.md>).