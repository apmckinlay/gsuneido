// Copyright (C) 2009 Suneido Software Corp. All rights reserved worldwide.
// requires On_ methods be forwarded to Scintilla as LibView does
ScintillaAddon
	{
	styleLevel: 49
	Init()
		{
		.marker_flag = .MarkerIdx(level: .styleLevel)
		}
	Styling()
		{
		return [[level: .styleLevel, marker: [SC.MARK_ROUNDRECT, back: CLR.BLUE]]]
		}
	ContextMenu()
		{
		return #('&Flag\tCtrl+F2', 'Next Flag\tCtrl+Shift+F2', 'Previous Flag\tShift+F2')
		}
	On_Flag()
		{
		lineno = .LineFromPosition()
		if .isFlag?(lineno)
			.MarkerDelete(lineno, .marker_flag)
		else
			.MarkerAdd(lineno, .marker_flag)
		.MarkersChanged()
		}
	isFlag?(lineno)
		{
		state = .MarkerGet(lineno)
		return (state & (1 << .marker_flag)) isnt 0
		}
	On_Next_Flag()
		{
		.locateFlag("MarkerNext")
		}
	On_Previous_Flag()
		{
		.locateFlag("MarkerPrevious")
		}
	locateFlag(method = false)
		{
		lineno = .LineFromPosition()
		count = 0
		startIndex = method is "MarkerPrevious" ? 29999 : 0 /*= from end */
		while (count++ < 700) /*= number of markers to attempt to check */
			{
			start = method is "MarkerPrevious" ? lineno - 1 : lineno + 1
			nextLine = this[method](start, 0xffff)
			if nextLine < 0
				nextLine = this[method](startIndex, 0xffff)

			if .isFlag?(nextLine)
				{
				.EnsureVisible(nextLine)
				.GotoLine(nextLine)
				return
				}
			lineno = nextLine
			}
		}
	}