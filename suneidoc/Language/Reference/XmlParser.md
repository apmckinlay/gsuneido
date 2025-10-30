<div style="float:right"><span class="toplinks"><a href="XmlParser/Methods">Methods</a></span></div>

### XmlParser

``` suneido
(text) => xmlNode
```

Uses the SAX style [XmlReader](<XmlReader.md>) to parse the text into a tree of [XmlNode](<XmlNode.md>)'s which then allow access to the contents.

For example:

``` suneido
text = '<html>
    <head>
    <title>Test</title>
    </head>
    <body>
    <img src="pic.jpg" align="right" />
    <p>some stuff</p>
    <p>some <b>more</b> stuff</p>
    </body>
    </html>'
xml = XmlParser(text)

xml.Name()
    => "html"

xml.head.title.Text()
    => "Test"

xml.body.p
    =>  <p>
            some stuff
        </p>
        <p>
            some 
            <b>
                more
            </b>
             stuff
        </p>

xml.body.p[1].Text()
    => "some more stuff"
    
xml.body.img[0].Attributes()
    => #(align: "right", src: "pic.jpg")

xml.body.img[0]._align
    => "right"
```