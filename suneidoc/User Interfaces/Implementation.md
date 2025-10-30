## Implementation

|     |
| --- |
| [Construction](<Implementation/Construction.md>) |
| [Size and Stretch](<Implementation/Size and Stretch.md>) |
| [Containment Requirements](<Implementation/Containment Requirements.md>) |
| [Writing a Hwnd Control](<Implementation/Writing a Hwnd Control.md>) |
| [Writing a WndProc Control](<Implementation/Writing a WndProc Control.md>) |
| [Writing Wrapper Controls](<Implementation/Writing Wrapper Controls.md>) |
| [Writing a Layout Control](<Implementation/Writing a Layout Control.md>) |
| [Writing a Controller](<Implementation/Writing a Controller.md>) |



-	controls that encapsulate a Windows control are derived from the Hwnd class (e.g. Button)
-	controls that act as a window procedure for a window are derived from the WndProc class (e.g. Scrollable, EtchedLine)
-	controls that manage other controls are derived from the Controller class (e.g. WorkSpace, Library View)
-	controls that are simply layout managers are derived from the Control class (e.g. Vert)
-	Window is a top level window - it should not be derived from