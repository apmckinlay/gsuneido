<div style="float:right"><span class="builtin">Builtin</span></div>

#### Image

``` suneido
( string ) => image
```

Create an instance of the specified image, 
which may be either a filename, or "book%name" or the actual image data.
Uses OleLoadPicture, which handles bmp, gif, jpg, ico, emf, and wmf images.

See also: 
[ImageControl](<../../../User Interfaces/Reference/ImageControl.md>),
[ImageFormat](<../../../Reports/Reference/ImageFormat.md>)