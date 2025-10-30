<div style="float:right"><span class="toplinks"><a href="/suneidoc/User Interfaces/Reference/WndProc/Methods">Methods</a></span></div>

### WndProc

Abstract base class for window procedure classes, derived from [Hwnd](<Hwnd.md>). A WndProc object acts as a window procedure via its Call method.

WndProc performs the following functions:

-	Translate WM Windows messages (numbers) to methods via Call.
-	Redirect notifications to the originating controls via COMMAND and NOTIFY.