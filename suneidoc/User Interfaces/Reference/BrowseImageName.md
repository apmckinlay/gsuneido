### BrowseImageName

``` suneido
(title = "", filter = "", file= false, showimage= true, opendirmsg = "")
```

Open a dialog window to can select an image file and display a preview of the selected image. Can also open another dialog window to change the current directory.

Returns the selected image file path or false if Cancel button is chosen.
`title`
: will display the contained text in the dialog title

`filter`
: Allows filtering the images files displayed in the dialog. The default filter is:  "bmp, gif, jpg, jpe, jpeg, ico, emf, wmf, tif, tiff, png, pdf". e.g. to set a filter to only JPEG images, set the string to "jpg, jpeg" (the space between is important)

`file`
: may contain: 
    1) false: to select the computer resources item in open directory dialog window
    2) true: to select the current directory in open directory dialog window
    3) text path: to select the passed directory path in open directory dialog window
    The text path may be:
    a) a full file path and a file name 
    b) a full file path to select a directory 
    c) a file name
    d) nothing

`showimage`
: true/false to display or not, the image in the dialog

`opendirmsg`
: a message displayed in the open directory dialog window, not the title dialog window