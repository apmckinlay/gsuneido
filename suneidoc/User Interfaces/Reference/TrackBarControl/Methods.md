### [TrackBarControl](<../TrackBarControl.md>) - Methods
`SetRange(min, max, redraw = true)`
: sets the range (minimum and maximum positions) for the slider. If redraw true, the slider is redrawn after the range is set. Default is 0 and 100

`Set(position)`
: sets the current position of the slider

`Get()`
: return the current position of the slider

`ClearTics(redraw = true)`
: removes the current tick marks from a trackbar

`SetTic(position)`
: sets the position of a tick mark. Returns true if the tick mark is set or false otherwise.

`SetTicFreq(freq, position)`
: sets the interval frequency for tick marks in a trackbar. For example, if the frequency is set to two, a tick mark is displayed for every other increment in the trackbar's range. The default setting for the frequency is one; that is, every increment in the range is associated with a tick mark.

`ClearSel(redraw = true)`
: clears the current selection in a trackbar.

`SetSel(start, end, redraw = true)`
: sets the starting and ending positions for the current selection

`GetSelStart()`
: sets the starting position of the current selection

`GetSelEnd()`
: sets the ending position of the current selection