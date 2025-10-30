## Formats

Report output is done through format items.  These format items range from drawing a line or printing some text, to producing a complete report from a database query.  Simple, low level formats use the Windows API to output to a device context.  This device context may be the screen or the printer depending on where the report is being sent.  More complicated formats are often built up from combinations of simpler formats.

The basic layout building blocks for reports are:

`Horz, Vert` - the basic containers for laying out multiple formats

`Hfill, Vfill` - leave stretchable blank space

`Hskip, Vskip` - leave fixed size blank space 

The low level formats you can use include:

`Number` - to output numbers using a specified mask

`Text` - to output a single line of text

`Wrap` - to output long strings of text wrapped onto multiple lines

`ShortDate` - to output a date (format determined by your Windows settings)

`Hline, Vline` - to draw horizontal and vertical lines

`Rect` - to draw a rectangle around another format item

The more complex formats include:

`Hfields, Vfields` - output a list of fields labeled with their prompts

`Query` - display a database query, normally using ColHeads and Row

`PageHead` - standard page header

`ColHeads` - standard column heads

`Row` - print fields in columns

`Total` - print an underlined total

Format specifications are written as objects, for example:

``` suneido
#( Text, "hello world" )
```

The first value in a specification object is the name of the format. If there are no arguments to the format, you can simply write the name of the format. i.e. #(Vskip) can be written as 'Vskip'.

Note: The actually format names in the library have a suffix of Format, e.g. the Text format is found in the library as TextFormat.

Special Items

|  |  | 
| :---- | :---- |
| `"pg"` | start a new page | 
| `"pg0"` | start a new page and reset the page number (the next page will be 1) | 
| `"pge"` | start a new even numbered page | 
| `"pgo"` | start a new odd numbered page | 


pge and pgo will generate a blank page if necessary.