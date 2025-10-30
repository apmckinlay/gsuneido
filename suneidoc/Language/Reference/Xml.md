### Xml

``` suneido
(tag, content = false [, attribute: value ...]) => string
```

Create XML (or XHTML or HTML) strings.

If there is no content, an *empty element* is returned e.g. \<tag />

For example:

``` suneido
Xml("br")
    => <br />'

Xml("p", "hello world")
    => '<p>hello world</p>'

Xml("p", "hello world", font: "serif", size: 10)
    => '<p font="serif" size="10">hello world</p>'
```

Calls to Xml can be nested:

``` suneido
Xml('b', Xml('i', "hello"))
    => '<b><i>hello</i>&lt/b>'
```

Use a tag starting with "?" for processing instructions:

``` suneido
Xml("?works" data: "test.wks")
    => ''
```

Use a tag of "!--" for comments:

``` suneido
Xml("!--", "this is a comment")
    => ""
```

Note: Xml does <u>not</u> encode & \< > " as character entities.

See also:
[Using Xml to Create XML and XHTML](<../../Cookbook/Using Xml to Create XML and XHTML.md>)