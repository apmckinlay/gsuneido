// Copyright (C) 2012 Suneido Software Corp. All rights reserved worldwide.
EnhancedButtonControl
	{
	WindowClass: 'SuBtnfaceArrowNoDblClks'

	LBUTTONDOWN()
		{
		super.LBUTTONDOWN()
		.SetPressed?(true)
		.Send(.GetCommand())
		// the button could be destroyed/rebuilt by the command, e.g. Vert.Insert
		if .Member?('Hwnd')
			TrackMouseEvent(Object(cbSize: TRACKMOUSEEVENT.Size(),
				dwFlags: TME.HOVER, hwndTrack: .Hwnd))
		return 0
		}

	LBUTTONUP()
		{
		.SetPressed?(false)
		.Send(.GetCommand() $ '_MouseUp')
		.Repaint()
		return 0
		}

	MOUSELEAVE()
		{
		.LBUTTONUP()
		return super.MOUSELEAVE()
		}

	MOUSEHOVER()
		{
		if .GetPressed?()
			.Send(.GetCommand() $ '_MouseHold')
		return 0
		}
	}
