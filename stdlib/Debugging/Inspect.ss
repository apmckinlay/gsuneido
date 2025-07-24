// Copyright (C) 2000 Suneido Software Corp. All rights reserved worldwide.
// Uses InspectControl
// defaults to dialog, use Inspect.Window(...) for a non-modal window
PassthruController
	{
	CallClass(x, title = "Inspect", hwnd = 0)
		{
		ToolDialog(hwnd, Object(this, x, :title, :hwnd), border: 0)
		}
	Window(x, title = "Inspect", hwnd = 0)
		{
		Window(Object(this, x, :title, :hwnd),
			keep_placement: 'Inspect', exStyle: WS_EX.TOPMOST)
		}
	New(x, title, hwnd)
		{
		super(Object("Inspect", x, title, hwnd))
		}
	}