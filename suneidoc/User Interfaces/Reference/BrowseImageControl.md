### BrowseImageControl

``` suneido
(title = '', filter = "", width = 30, file = false, showimage = true,
    status = "", opendirmsg = "")
```

Use BrowseImageControl when you want the user to enter an existing image file name with the option of browsing for it in a dialog that display a preview of the selected image.
`title`
: will display the contained text in the dialog title

`filter`
: Allows filtering the images files displayed in the dialog. The default filter is:  "bmp, gif, jpg, jpe, jpeg, ico, emf, wmf, tif, tiff, png, pdf". e.g. to set a filter to only JPEG images, set the string to "jpg, jpeg" (the space between is important)

`width`
: is multiplied by the average character width

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
: whether or not to display the image in the dialog

`status`
: allows you to set a string to display as a tip or to appear in the status bar when the field has the focus.

`opendirmsg`
: is a message displayed in the open directory dialog window, not the title dialog window

The browse button uses [BrowseImageName](<BrowseImageName.md>) to bring up the browse image dialog

Methods:
`Set(file)`
: to set a file value in the field text

`SetFilter(filter)`
: to set the filter for the images that will be displayed in the dialog

`SetStatus(status)`
: to set a tip for the field

`Get()`
: to retrieve the field text value

`GetFileName()`
: to retrieve only the file name from the field text

`GetFilePath()`
: to retrieve only the file path from the field text

`Dirty?(dirty = "")`
: return true when the data is modified

`SetReadOnly(boolean)`
: to set the control to read only or not