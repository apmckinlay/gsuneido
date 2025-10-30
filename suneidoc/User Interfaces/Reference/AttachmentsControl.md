### AttachmentsControl

``` suneido
(control)
```

A PassthruController that implements ImageDropFileList to handle multiple files dropped onto one of its child OpenImageControls. Files are placed into sequential child OpenImageControl's starting with the one that was dropped on.

For example:

``` suneido
(Attachments (Horz OpenImage OpenImage OpenImage))
```

See also: 
[ImageFormat](<../../Reports/Reference/ImageFormat.md>),
[OpenImageControl](<OpenImageControl.md>)