### [MimeBase](<../MimeBase>) - Methods
`Content_Transfer_Encoding(value) => this`
: 

`To(string) => this`
: 

`From(string) => this`
: 

`Date(date = Date()) => this`
: 

`Subject(string) => this`
: 

`Message_ID( [ string ] ) => this`
: The above methods set the corresponding header field. Date and Message_ID have default values.

`SetPayload(string) => this`
: 

`AddHeader(name, string, [extra: string]) => this`
: For example: `AddHeader("Content-Disposition", "attachment")`   
to get: `Content-Disposition: attachment`   
or: `AddHeader("Content-Disposition", "attachment", filename: "test.txt")`   
to get: `Content-Disposition: attachment; filename="test.txt"`

`AddExtra(extra: string) => this`
: For example: `AddExtra("Content-Type", charset: "us-ascii")`   
would result in `Content-Type: text/plain; charset="us-ascii"`

`Base64() => this`
: Sets the encoding to Base64 and adds Content-Transfer-Encoding: base64

`ToSring() => string`
: Returns the resulting message as a string.

Most of the methods return this so they can be "strung" together as in:

``` suneido
MimeBase("application", "octet-stream").
    AddHeader("Content-Disposition", "attachment", filename: "test.txt").
    SetPayload("second part").Base64().ToString()
```