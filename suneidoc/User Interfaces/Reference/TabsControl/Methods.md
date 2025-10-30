### [TabsControl](<../TabsControl.md>) - Methods
`ConstructTab(i)`
: Creates the control for the tab at position *i*.  The control is not visible until the tab is selected. If the tab's control has already been constructed, ConstructTab does nothing.

`Constructed?(i)`
: Returns true if the tab's control is constructed, false otherwise.

`GetControl(i = false)`
: Returns the control for the tab at position *i*.  If i is false, then the control for the currently selected tab is returned.

`GetSelected()`
: Returns the index of the selected tab (starting at 0).

`Select(i)`
: Makes the current tab the tab at index *i* and displays its contents.

`DestroyTab(i)`
: &nbsp;

`SetEnabled(bool)`
: enable/disable tabs

`SetReadOnly(bool)`
: if true the controls in the tab are readonly

`SetVisible(bool)`
: if false the tabs are not visible

`GetTabCount()`
: get the number of tabs in the TabsControl

`SetImageList(images)`
: set the image list

`SetImage(i, img)`
: set an image in a tab