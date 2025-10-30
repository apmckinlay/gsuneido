### [XmlBuilder](<../XmlBuilder>) - Methods
`_(string) => this`
: Adds the string to the output. Useful for mixed content:

``` suneido
XmlBuilder(indent: 4).div
    {
    .h1("Heading")
    ._("stuff")
    .p("paragraph")
    }
=>  <div>
        <h1>Heading</h1>
        stuff
        <p>paragraph</p>
    </div>
```

`Comment(string) => this`
: Adds a comment to the output:

``` suneido
XmlBuilder().Comment("hello world")
    => "<!-- hello world -->"
```

`Declare(@args) => this`
: If no arguments are supplied a standard XHTML declaration is added:

``` suneido
XmlBuilder().Declare()
    =>  <!DOCTYPE html PUBLIC "-//W3C//DTD XHTML 1.0 Transitional//EN"
            "http://www.w3.org/TR/xhtml1/DTD/xhtml1-transitional.dtd">
```

&#8203;
: Otherwise, the first argument is the tag name and any remaining arguments are added. String arguments are quoted, symbol arguments are not. For example, this would produce the default:

``` suneido
XmlBuilder().Declare('DOCTYPE', #html, #PUBLIC,
    "-//W3C//DTD XHTML 1.0 Transitional//EN",
    "http://www.w3.org/TR/xhtml1/DTD/xhtml1-transitional.dtd")
```

`Instruct(instruction, attribute: value ...) => this`
: The first un-named argument is the instruction. The default is "xml", in which case the attribute defaults are version: "1.0" and encoding: "US-ASCII". Named arguments are taken as attributes. Attribute values are XmlEntityEncode'd

``` suneido
XmlBuilder().Instruct()
    => "<?xml encoding="US-ASCII" version="1.0">"

XmlBuilder().Instruct('XML', version: '1.1', encoding: 'UTF-8'
    => "<?XML encoding="UTF-8" version="1.1">"
```

`Default(@args)`
: This is the method called when you call xmlBuilder.tag It can be called directly if the tag is not a valid method name. For example:

``` suneido
XmlBuilder().Default("SOAP:Encoding") { "..." }
    => "<SOAP:Encoding>...</SOAP:Encoding>"
```

&#8203;
: The first un-named argument is the tag. If there is a block argument, it is called (with object.Eval2). If the block returns a value (other than this) it will be added to the output. If there is no block a second un-named argument is taken as the text content. 

`ToString() => string`
: Returns the output string.