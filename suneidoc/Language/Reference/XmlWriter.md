### XmlWriter

An [XmlReader](<XmlReader.md>) content handler that builds an XML string from the callbacks.

Empty elements, e.g. StartElement("br"), EndElement("br") with no characters between are converted to a single tag e.g. "\<br />

The resulting string is retrieved with the GetText method.

For example:

``` suneido
xr = new XmlReader
xw = new XmlWriter
xr.SetContentHandler(xw)
xr.Parse('<p>Hello <font size="3">new</font> world</p>')
return xw.GetText()

=> '<p>Hello <font size="3">new</font> world</p>'
```

Encodes \&amp; \&lt; \&gt; \&quot; character entities.

Note: Currently, no newlines or indenting are added - the result is one long line.

See also: [Xml](<Xml.md>)