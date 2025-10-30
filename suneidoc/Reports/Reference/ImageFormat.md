### ImageFormat

``` suneido
(image = false, scale = 1, width = false, height = false,
    stretch = false, background = false)
```

Uses [Image](<../../Language/Reference/Image.md>)
to output the specified image, which may be either a filename, or the actual image data.
Uses OleLoadPicture, which handles bmp, gif, jpg, ico, emf, and wmf images.

If you specify one of **width** or **height**, the other dimension will be sized to maintain the proportions. If you specify both **width** and **height** the image will be sized to fit within those dimensions, maintaining proportions unless  **stretch** is true.

**background** can be specified as an object containing the x and y position (in inches) that the image should be printed at on the page with the specified **width** and **height**. For example:

``` suneido
(Image 't4.gif',
    background: (x: .25, y: .25), width: 8, height: 10)
```

See also: 
[ImageControl](<../../User Interfaces/Reference/ImageControl.md>)