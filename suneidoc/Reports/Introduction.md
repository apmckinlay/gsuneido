## Introduction

Reports layout is done using the "Boxes and Stretch" pattern .

Reports are normally run through Params.  The basic usage is:

``` suneido
Params( item, ... )
```

An item may be a printable format like Text, a container like Vert or Horz, or a generator like Query.

A very simple report would be:

``` suneido
Params( #( Text, "Hello world!" ) )
```

which would simply print "Hello world!" at the top left of the page.

If a title is provided, a default page header including the date and time and the page number will be produced:

``` suneido
Params( #( Text, "Hello world!" ), title: "Test" )
```

or you can provide your own format for the page header:

``` suneido
Params( #( Text, "Hello world!" ), header: #( Text, "Test" ) )
```

Note: Titles have been omitted from the following examples to keep them short.