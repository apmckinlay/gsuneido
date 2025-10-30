### ImageControl

``` suneido
(image = "", message = "no image", stretch = false, allowDrop = false)
```

Uses [Image](<../../Language/Reference/Image.md>) to display the specified image.
`image`
: May be either a filename, or the actual image data, or a string of the form "book%name"

`message`
: If there is no image ("") or the image is invalid then the **message** is displayed instead.

`stretch`
: If **stretch** is true, then the image will be stretched to fit the control. Otherwise, the image proportions will be maintained and the image will be centered in the control.

`allowDrop`
: If true, ImageControl will call DragAcceptFiles on itself and respond to WM_DROPFILES by Send'ing 'ImageDropFiles' with wParam (hDrop). The controller must implement the drop behavior.

If xmin and ymin are <u>not</u> specified, the control is set to the same size as the image. If only one of xmin or ymin are specified, the other is set to maintain the proportions of the image.

For example:

``` suneido
ImageControl("c:/images/face.gif")
```

might produce:

![](<../../res/Image.gif>)

See also: 
[ImageFormat](<../../Reports/Reference/ImageFormat.md>),
[OpenImageControl](<OpenImageControl.md>)