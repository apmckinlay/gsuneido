### [StatusControl](<../StatusControl.md>) - Methods
`Get() => string`
: Returns text of the status.

`GetValid() => boolean`
: Returns false if status is invalid (red color), else returns true.

`Set(text)`
: Sets the displayed text in the statusbar.

`SetValid(valid)`
: Sets status invalid (red color) if *valid* argument is false, if true set status valid (COLOR.BTNFACE).