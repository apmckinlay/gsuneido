### AtLeastOneControl

``` suneido
(control)
```

A [PassthruController](<PassthruController.md>) that manages several true/false controls (e.g. [CheckBoxControl](<CheckBoxControl.md>) or [ButtonToggleControl](<ButtonToggleControl.md>)) and prevents changing all the controls to false.

If there are exactly two controls, toggling the only true control will automatically set the other control to true.

If there are more than two controls, then you cannot toggle to only true control. (You must set another control to true first.)

**Note**: This does <u>not</u> guarantee that you always have at least one control set to true. For example, when the controls are first created they may be all false.