### Format

Abstract base class for formats.

Derived classes define:
`Generator?( ) => true or false`
: Should return true if the format is a generator, false otherwise. The default definition returns false. Normally, generators should be derive from Generator which defines Generator? for you.

`GetSize(data = false) => #(w: , h: , d: )`
: Should return the width, height, and descent of the format item.

`OnPage()`
: Called by Report when a format item is fitted onto a page. Note: page related processing should not be done in format.New because the item may not fit on the current page. And it is often too late to do this type of processing in Print.

`Print(x, y, w, h, data = false)`
: Do the actual output.

`Variable?( ) => true or false`
: Should return true if the format's size depends on it's data, false if the size is fixed. The default definition returns false.

`Xstretch`
: If Xstretch is greater than 0 the format is stretchable. If a Horz is printed with a larger width than it requires, the extra space will be distributed among its stretchable items, proportional to their Xstretch values.

`Ystretch`
: If Ystretch is greater than 0 the format is stretchable. If a Vert is printed with a larger height than it requires, the extra space will be distributed among its stretchable items, proportional to their Ystretch values.