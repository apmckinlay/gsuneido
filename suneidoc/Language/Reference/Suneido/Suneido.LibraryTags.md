<div style="float:right"><span class="builtin">Builtin</span></div>

#### Suneido.LibraryTags

``` suneido
(@args)
```

For example:

``` suneido
Suneido.LibraryTags("trial", "sujs")
```

Suneido.LibraryTags sets the current library tags which affect how code is loaded from libraries. With no arguments it resets the library tags to their initial startup state.

Code is loaded from libraries in two stages. When running client-server the first stage happens on the server and the second stage happens on the client. But the two stages happen even when running standalone. The first stage gets a list of library records for a given global name from the libraries in use, with the current set of tags.

For example, loading "Foo" when [Libraries()](<../Libraries.md>) are alib, blib and tags are xtag, ytag the maximum number of records that could possibly be returned are:

``` suneido
alib:Foo
alib:Foo__xtag
alib:Foo__ytag
blib:Foo
blib:Foo__xtag
blib:Foo__ytag
```

The second stage compiles this list of records. Later libraries/tags override earlier ones but can also reference the previous definition. See [Libraries](<../../Libraries.md>)

When running client-server, the tags set on the server will determine what code is sent to the client. However, a client can set different tags with the result being the intersection of the tags. e.g. if the server has a,b,c and the client has b,c,d the result will be the equivalent of b,c.

The current library tags can be retrieved with [Suneido.Info](<Suneido.Info.md>)("library.tags")

Note: LibraryTags calls [Unload()](<../Unload.md>)