### Composites

Output formats actually make Windows API calls to display their contents.  Composite formats use other formats to do the actual output.  

For example, here is a simple page heading format:

``` suneido
HorzFormat
    {
    Xstretch: 1
    New( )
        {
        super(@.format( ));
        }
    format(title)
        {
        return Object(
                Object( 'Text', DateString( ) ),
                'Hfill',
                Object('Text', "Page " $ _report.GetPage( ) )
                );
        }
    }
```

Note: In this example the body of the format method 
could have been placed directly in New's super call.
However, this is not possible if more complex construction is required.

Of course, hybrids are possible where some of the output is delegated to other formats 
and some is done directly.
For example, ColHeads uses Text formats to output the headings, 
but draws the lines directly.