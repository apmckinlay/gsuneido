### OpenFileName

``` suneido
(filter = "", hwnd = false, flags = false, multi = false, 
    title = "Open", file = "", initialDir = "")
```

Opens the standard Windows Open file dialog.
`filter`
: allows filtering the file names displayed in the dialog, the default filter is `"All Files (*.*)\000*.*\000"`

`flags`
: sets the Windows Open file dialog options, the default is `OFN.FILEMUSTEXIST | OFN.PATHMUSTEXIST | OFN.HIDEREADONLY | OFN.NOCHANGEDIR`

`multi`
: If set to true, the user is allowed to select multiple files.

`title`
: will display in the dialog title

`file`
: a file name

`initialDir`
: a directory path to select as default opening the dialog

Used by [OpenFileControl](<OpenFileControl.md>).

See also:
[SaveFileName](<SaveFileName>)