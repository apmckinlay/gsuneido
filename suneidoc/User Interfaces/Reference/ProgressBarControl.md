<div style="float:right"><span class="toplinks"><a href="/suneidoc/User Interfaces/Reference/ProgressBarControl/Methods">Methods</a></span></div>

### ProgressBarControl

``` suneido
(style)
```

ProgressBarControl is a simple wrapper for the Windows Progress Bar control.  A style parameter can be specified to add to the default WS.VISIBLE style. Two progress bar styles that can be applied are PBS.SMOOTH and PBS.VERTICAL.

Here's a very simple example:

``` suneido
pb = ProgressBarControl().Ctrl
for (i = 0; i < 9; ++i)
    {
    pb.StepIt()
    Sleep(100)
    }
```