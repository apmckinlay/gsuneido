<div style="float:right"><span class="builtin">Builtin</span></div>

#### string.Unescape

``` suneido
() => string
```

Returns a copy of the string with *escape* sequences 
replaced with their corresponding characters.

-	\xhh for hex character values, where h is 0-9, a-f, or A-F
-	\t for tab, \n for newline, \r for carriage return
-	\\' for single quote, \\" for double quote
-	\0 for a nul (zero) character
-	\\\\ for a literal backslash


This is the same translation that Suneido's compiling does on string literals.