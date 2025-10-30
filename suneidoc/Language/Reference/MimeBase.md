<div style="float:right"><span class="toplinks"><a href="MimeBase/Methods">Methods</a></span></div>

### MimeBase

``` suneido
(maintype = "text", subtype = "plain")
```

Basic support for generating MIME messages.

MimeText and MimeMultiPart are derived (inherit) from MimeBase.

For example:

``` suneido
MimeBase("application", "octet-stream").
    AddHeader("Content-Disposition", "attachment", filename: "test.txt").
    SetPayload("second part").Base64().ToString()
```

would produce:

``` suneido
Content-Type: application/octet-stream
MIME-Version: 1.0
Content-Disposition: attachment; filename="test.txt"
Content-Transfer-Encoding: base64

c2Vjb25kIHBhcnQ=
```