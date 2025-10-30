### CodeFormat

``` suneido
( text, width = False, w = 2880, font = False )
```

Like `WrapFormat`,
except that wrapped lines are indented one more tab stop than the original.

width is in characters, based on average character width, 
and therefore dependent on the font.
w is in twips.  Specify one or the other.

Used by [LibraryFormat](<LibraryFormat.md>).