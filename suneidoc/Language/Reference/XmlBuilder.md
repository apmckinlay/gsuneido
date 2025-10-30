<div style="float:right"><span class="toplinks"><a href="XmlBuilder/Methods">Methods</a></span></div>

### XmlBuilder

``` suneido
(indent = 0, margin: 0)
```

Used to construct XML strings. For example:

``` suneido
XmlBuilder(indent: 4).
    Instruct().
    Declare().
    Comment('test').
    html
        {
        .head { .title { "My Title" } }
        .body
            { 
            .img(src: "pic.jpg")
            .h1("First")
            .p { "The first paragraph" }
            .h1 { "Second" }
            .p { "A second paragraph" }
            .S("stuff")
            }
        }

=>  <?xml encoding="US-ASCII" version="1.0">
    <!DOCTYPE html PUBLIC "-//W3C//DTD XHTML 1.0 Transitional//EN" "http://www.w3.org/TR/xhtml1/DTD/xhtml1-transitional.dtd">
    <!-- test -->
    <html>
        <head>
            <title>My Title</title>
        </head>
        <body>
            <img src="pic.jpg" />
            <h1>First</h1>
            <p>The first paragraph</p>
            <h1>Second</h1>
            <p>A second paragraph</p>
            stuff
        </body>
    </html>
```

tags are written as method calls on the builder. Attributes are written as named arguments. Content can be supplied as a first un-named argument e.g. .h1("First") or as a block returning a string e.g. .h1 { "Second" } Raw text can be added with .S

If indent is not zero the output is formatted, otherwise no indenting or newlines will be added. For example:

``` suneido
XmlBuilder().div { .h1("Hello") }
    => <div><h1>Hello</h1></div>

XmlBuilder(indent: 4).div { .h1("Hello") }
    => <div>
           <h1>Hello</h1>
       </div>
```

A margin may be specified to indent the entire output.