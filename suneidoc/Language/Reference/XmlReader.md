### XmlReader

A simple implementation of a SAX (Simple API for XML) XML reader which parses XML strings and does callbacks on a content handler.
`SetContentHandler(handler)`
: Sets the content handler. If this method is not called, the default content handler is 
[XmlContentHandler](<XmlContentHandler.md>) which simply ignores callbacks.

`Parse(text)`
: Parses the XML text and makes callbacks to the content handler of StartElement(qname, atts), EndElement(qname), Characters(string), and IgnorableWhitespac(string).

For example:

``` suneido
handler = XmlContentHandler
    {
    StartElement(qname, atts)
        { Print("START", qname, atts) }
    EndElement(qname)
        { Print("END", qname) }
    Characters(s)
        { Print("CHARS", Display(s)) }
    }
xr = new XmlReader
xr.SetContentHandler(handler)
xr.Parse('<p>Hello <font size="3">new</font> world</p>')

=>  START p #()
    CHARS "Hello "
    START font #(size: "3")
    CHARS "new"
    END font
    CHARS " world"
    END p
```

`<? ... ?>` processing instructions; and `<!-- ... -->` comments are ignored.

The contents of `<![CDATA[ ... ]]>` sections are returned as Characters (with no other processing).

Used by [XmlRpc](<XmlRpc.md>) and [XmlParser](<XmlParser.md>).