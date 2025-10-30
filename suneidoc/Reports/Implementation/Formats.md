### Formats

Formats are constructed using _report.Construct()

Formats use the following methods:

``` suneido
GetSize(data = ?) => Object(w:, h:, d:)
```

Returns the width, height, and descent of the format.

``` suneido
Print(x, y, w, h, data = ?)
```

Prints the format at position x,y with the specified width and height.  If the format is stretchable, the width and height may be larger than what GetSize returned.  

``` suneido
Variable?( ) => true or false (default false)
```

Returns true if the format's size depends on the data (e.g. Wrap).

``` suneido
Generator?( ) => true or false
```

Returns true if the format is a generator.

and the following public members (with their default values):

``` suneido
X = 0
Y = 0
Xstretch = 0
Ystretch = 0
Data
Field
```

Report.Construct automatically copies certain members from the format specification to the constructed format instance:

``` suneido
x
y
xstretch
ystretch
field
```

Formats have access to their containing report via:

``` suneido
_report
```

For example:

``` suneido
_report.GetDC( ) // get the current device context

_report.Construct( ) // construct a format

_report.SelectFont( )

_report.GetFont( )
```

If a format changes the current font (using _report.SelectFont) 
it is also responsible for changing it back:

``` suneido
oldfont = _report.SelectFont( .font );
...
_report.SelectFont( oldfont );
```

Formats can be very simple, for example here is Vfill:

``` suneido
Format
    {
    Ystretch: 1
    }
```

Vskip is a little more complicated because it allows you to override the height:

``` suneido
Format
    {
    Size: #( h: 240, w: 0, d: 0 );
    New( height = false )
        {
        if ( height isnt false )
            .Size = Object( h: height * 1440, w: 0, d: 0 );
        }
    GetSize( data = false )
        {
        return .Size;
        }
    }
```

The general outline for an actual printing format should be:

``` suneido
Format
    {
    New(data = false, ...)
        {
        .Data = data;
        ...
        }
    GetSize(data = "")
        {
        if (.Data isnt false)
            data = .Data;
        ...
        }
    Print(x, y, w, h, data = "")
        {
        if (.Data isnt false)
            data = .Data;
        ...
        }
    }
```

Notice how the data may be supplied either when the format is created, or when it is used.
This allows the Flyweight pattern to be used, so the formats can be reused.

For example here is a simple text format:

``` suneido
DrawTextFormat
Format
    {
    New(data = false, font = false)
        {
        .Data = data;
        .font = font;
        }
    GetSize(data = "")
        {
        if (.Data isnt false)
            data = .Data;
        data $= "";
        oldfont = _report.SelectFont(.font);
        dc = _report.GetDC();
        flags = DT.SINGLELINE +
            DT.CALCRECT + DT.EXTERNALLEADING;
        DrawTextEx(dc, data, -1, r = Object(), flags, 0);
        _report.SelectFont(oldfont);
        return Object(h: r.bottom, w: r.right, d: tm.Descent);
        }
    Print(x, y, w, h, data = "")
        {
        if (.Data isnt false)
            data = .Data;
        data $= "";
        oldfont = _report.SelectFont(.font);
        flags = DT.SINGLELINE;
        rect = Object(left: x, right: x + w, top: y, bottom: y + h);
        DrawTextEx(_report.GetDC(), data, -1, rect, flags, 0);
        _report.SelectFont(oldfont);
        }
    }
```

This could be used like:

``` suneido
( DrawText "hello world" font: ( name: "Arial", size: 20 ) )
```