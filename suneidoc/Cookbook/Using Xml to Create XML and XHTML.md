## Using Xml to Create XML and XHTML

**Category:** Coding

**Problem**

You need to create HTML (or XML or XHTML) strings.

**Ingredients**

[Xml](<../Language/Reference/Xml.md>) function, string concatenation, blocks, [QueryAccum](<../Database/Reference/QueryAccum.md>)

**Recipe**

The Xml function in stdlib provides a simple yet flexible way to generate XML.

The simplest form is a tag with no content:

``` suneido
Xml("br")

    => "<br />"
```

The next step up is wrapping a string with tags:

``` suneido
Xml("em", "hello")

    => "<em>hello</em>"
```

You can also add attributes:

``` suneido
Xml("td", "hello world", valign: "top", bgcolor: "red")

    => '<td valign="top" bgcolor="red">hello world</td>'
```

Attributes can be used for inline CSS styles:

``` suneido
Xml("span", "hello world", style: "color: red; font-size: 20pt;")

    => '<span style="color: red; font-size: 20pt">hello world</span>'
```

Calls to Xml can of course be nested and concatenated:

``` suneido
Xml("table",
    Xml("tr", Xml("td", "top left") $ Xml("td", "top right")) $ "\n" $
    Xml("tr", Xml("td", "bottom left") $ Xml("td", "bottom right")),
    border: 1)

    => '<table border="1"><tr><td>top left</td><td>top right</td></tr>
       <tr><td>bottom left</td><td>bottom right</td></tr></table>'
```

Notice that the attributes must follow the content. This is because in Suneido function calls un-named arguments (i.e. the content) must come before named arguments.

Xml will also accept a *block* as the content. This allows you, for example, to wrap tag(s) around content produced by a loop.

``` suneido
Xml("table")
    {
    Xml('tr', Xml('td', "Table") $ Xml('td', "Size")) $ "\n" $
    QueryAccum("tables", "")
        { |s, x|
        s $ Xml('tr', Xml('td', x.tablename) $ Xml('td', x.totalsize)) $ "\n"
        }
    }
```

Notice the addition of newlines in the last two examples. Although this is not strictly necessary, it makes the resulting XML much more human readable.

**Discussion**

Using Xml instead of straight string manipulation has a number of advantages:

-	you don't have to type so many angle brackets!
-	it automatically generates matching end tags, avoiding missing or mismatched tags
-	empty elements (with no content) automatically generate the closing slash required by XML and XHTML (e.g. <br />
-	all attribute values are quoted as required by XML and XHTML
-	the code is more *readable*


**See Also**

[Building a Large String a Piece at a Time](<Building a Large String a Piece at a Time.md>)