### SaveFileName

``` suneido
(filter = "", hwnd = false, flags = false, 
    title = "Save", ext = "", file = "", initialDir = "")
```

Opens the standard Windows Save file dialog.
`filter`
: specify which file names are displayed in the dialog, the default is `"All Files (*.*)\000*.*\000"`

`flags`
: sets the Windows Save file dialog options, the default is `OFN.OVERWRITEPROMPT | OFN.PATHMUSTEXIST | OFN.HIDEREADONLY | OFN.NOCHANGEDIR`

`title`
: displayed in the dialog title

`ext`
: a default file extention

`file`
: a file name

`initialDir`
: initial directory path

Used by [SaveFileControl](<SaveFileControl.md>).

See also:
[OpenFileName](<OpenFileName>),