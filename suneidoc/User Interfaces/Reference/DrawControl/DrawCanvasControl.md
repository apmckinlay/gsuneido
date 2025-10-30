#### DrawCanvasControl

Used by DrawControl to process Canvas mouse events. Derived from [CanvasControl](<../CanvasControl.md>).

Receives LBUTTONDOWN, MOUSEMOVE, and LBUTTONUP messages, extracts x and y from 
lParam, and then passes them on to the current "tracker", which is set with the
SetTracker method. If the tracker MouseUp returns a canvas item, it is added to the 
canvas.