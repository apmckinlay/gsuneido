## Building a Large String a Piece at a Time

**Category:** Coding

**Problem**

You need to build a large string a piece at a time.

**Ingredients**

strings, string concatenation

**Recipe**

Simply concatenate the pieces together. For example:

``` suneido
s = ""
while (false isnt piece = NextPiece())
    s $= piece
```

Note: If you have the pieces in an object, consider using [object.Join](<../Language/Reference/Object/object.Join.md>)

**Discussion**

This recipe is almost too simple to include. The reason it's here, is that most other languages with immutable strings (e.g. Java, Python, Lua) specifically tell you <u>not</u> to do it this way. This is because, with a normal string implementation, every time you concatenate, a new string is created and the contents of the two strings is copied into it. In this situation, building a large string is very inefficient. Java recommends using a StringBuffer instead of normal strings. Python recommends creating a list of pieces and then concatenating them. A Lua technical notes suggests using a somewhat complicated stack of string pieces using a "tower of hanoi" system. All this just to build a string?

With Suneido you don't have to worry about it. Suneido's string implementation delays concatenation until it is necessary, building a linked list of pieces instead. This is completely transparent - done for you automatically.