### [MimeMultiPart](<../MimeMultiPart>) - Methods
`Attach(mime) => this`
: Add another part. The argument will normally be another Mime object, although it can be any object that with a suitable ToString method.

`AttachFile(filename) => this`
: Reads the file and then creates and attaches an appropriate Mime object (either MimeText or MimeBase) using MimeType to determine the type from the filename extension. If the filename extension is not recognized the type will be set to application/octet-stream. Binary files are base64 encoded.