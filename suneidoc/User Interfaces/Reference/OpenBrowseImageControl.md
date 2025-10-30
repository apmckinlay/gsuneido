### OpenBrowseImageControl

``` suneido
(filter = "", file = "", status = "", noembed = false, showimage = true, 
    opendirmsg = "")
```

Uses [ImageControl](<ImageControl.md>) to display an image and [BrowseImageControl](<BrowseImageControl.md>) to select it.

Double clicking brings up a dialog allowing browsing to selected an image file name with a preview image of the selected file name.

The image can be saved either as a file name (linked) or as the actual image data (embedded). Embedded images must be less than 64k.
`filter`
: Allows filtering the images files displayed in the dialog. The default filter is:  "bmp, gif, jpg, jpe, jpeg, ico, emf, wmf, tif, tiff, png, pdf". e.g. to set a filter to only JPEG images, set the string to "jpg, jpeg" (the space between is important)

`file`
: may contain: 1) a full file path and a file name 2) a full file path to select a directory 3) a file name 4) nothing

`status`
: allows you to set a text to display as a tip or to appear in the status bar when the field has the focus

`noembed`
: if true only allow linking images

`showimage`
: if true show a preview of the image in the dialog

`opendirmsg`
: to set a text to display in open directory dialog

Uses [ImageControl](<ImageControl.md>) to display an image and [BrowseImageControl](<BrowseImageControl.md>) to select it.

Double clicking brings up a dialog allowing browsing to selected an image file name with a preview image of the selected file name.

The image can be saved either as a file name (linked) or as the actual image data (embedded). Embedded images must be less than 64k.

Methods:
`Set(file)`
: to set a file value in field text

`SetFilter(filter)`
: to set a filter for the images files that will be displayed in the dialog

`SetStatus (status)`
: to set a tip for the field

`SetShowimage (shoimage)`
: true/false to show a preview of the image

`SetOpendirmsg (opendirmsg)`
: to set an open directory message

`Get()`
: to retrieve the field text

`GetFileName()`
: to retrieve only the file name from the field text

`GetFilePath()`
: to retrieve only the file path from the field text

`SetReadOnly(bool)`
: to set the control to read only or not

Example:

``` suneido
Controller
    {
    Title: 'test OpenBrowseImage'
    Xmin: 620   // window width
    Ymin: 240   // window height
    New(dialogtip = "choose second image")  // pass parameters from external
        {
        // values set at runtime
        .img1 = .Data.Vert.Form1.Img1         // control path
        .img1.SetFilter("jpg, bmp")           // set image types filter
        .img2 = .Data.Vert.Form1.Img2
        .img2.SetStatus(dialogtip)            // dialog tip

        .img2.Set('c:\\Tempfile\\test.jpg')   // set OpenImage Img1 field value

        // with RecordControl it can be used the control name to manage controls
        // otherwise shoul be used the control path, ex. .img2 = .Data.Vert.Form1.Img2
        // Data in the default name of RecordControl
        .Data.GetControl("Img1").SetStatus("choice the first image file")
        .Data.GetControl("Img3").SetStatus("another image file to choose")
        for field in Object("Fld1", "Fld2", "Fld3", "Fld4")
            .Data.GetControl(field).SetVisible(false)
        }
    Controls:
    (Record
    (Vert
        (Skip 10)
        (Form name: 'Form1'
            (Static 'only an initial path' group: 0 )
            (Static 'complete \path\file name' group: 1 )
            (Static 'filename' group: 2 )
            (Static 'nothing' group: 3 ) nl
            (OpenBrowseImage name: 'Img1', xmin: 150, ymin: 150, file: 'c:\\Tempfile\\', noembed: true, showimage: true, group: 0 )
            (OpenBrowseImage name: 'Img2', xmin: 150, ymin: 150, filter: 'jpg', showimage: false, opendirmsg: 'select best image' group: 1)
            (OpenBrowseImage name: 'Img3', xmin: 150, ymin: 150, file: 'test.jpg', showimage: true, filter: "jpg, jpeg, bmp, png, gif"
 group: 2 )
            (OpenBrowseImage name: 'Img4', status: 'last image to choose', noembed: true, xmin: 150, ymin: 150, group: 3 ) nl nl
            (Field name: 'Fld1', readonly:true, width: 15, group: 0 )
            (Field name: 'Fld2', readonly:true, width: 15, group: 1 )
            (Field name: 'Fld3', readonly:true, width: 15, group: 2 )
            (Field name: 'Fld4', readonly:true, width: 15, group: 3 )
        )
    )
    )
    //
    setField(field, value)
        {
         .Data.GetControl(field).Set(value)
         .Data.GetControl(field).SetVisible(true)
        }
    //
    Record_NewValue(fieldsource,value)
        {
        control = .Data.GetControl(fieldsource)
        switch control.Name
            {
        case "Img1":
            .setField("Fld1", value)
        case "Img2":
           {
            .setField("Fld2", control.Get())
            Print(.img2.GetFilePath())
            Print(control.GetFileName())
           }
        case "Img3":
            .setField("Fld3", value)
        case "Img4":
            .setField("Fld4", control.Get())
            }
        }
    }
```

See also: [OpenImageControl](<OpenImageControl.md>)