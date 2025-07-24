// Copyright (C) 2000 Suneido Software Corp. All rights reserved worldwide.
// base class for VertSplit and HorzSplit, not for users
Group
	{
	splitSaveName: false
	New(first, second, handle? = false, splitter = false, .splitSaveName = false,
		.splitSaveNameSuffix = ' - Split Position')
		{
		super(Object(first, .getSplitter(handle?, splitter), second))

		.ctrls = .GetChildren()
		if (.Dir is 'vert')
			{
			.ystretch0 = .ctrls[0].Ystretch
			.ystretch2 = .ctrls[2].Ystretch
			}
		else // horz
			{
			.xstretch0 = .ctrls[0].Xstretch
			.xstretch2 = .ctrls[2].Xstretch
			}
		.moveObservers = Object()
		}
	getSplitter(handle?, splitter)
		{
		if splitter isnt false
			return splitter
		return handle? ? "HandleSplitter" : "Splitter"
		}
	SetSplitSaveName(.splitSaveName, suffix = '')
		{
		if suffix isnt ''
			.splitSaveNameSuffix = suffix
		if .firstResize is true
			return false
		return .loadSplit()
		}
	UpdateSplitter(remove = false)
		{
		.Remove(1)
		.Insert(1, remove ? #Vert : 'Splitter')
		}

	firstResize: true
	Resize(x, y, w, h)
		{
		.r = Object(:x, :y, :w, :h)
		if .firstResize is true
			{
			.firstResize = false
			loaded = .loadSplit()
			if loaded is false and .Dir is "vert" and .ctrls[2].Base?(ScrollControl)
				.maximize_scrollable()
			}
		super.Resize(x, y, w, h)
		}
	loadSplit()
		{
		if .splitSaveName is false or .Destroyed?()
			return false
		prevSplit = UserSettings.Get(.splitSaveName $ .splitSaveNameSuffix, false)
		if prevSplit isnt false
			{
			.SetSplit(prevSplit)
			return true
			}
		return false
		}
	maximize_scrollable()
		{
		hwnd = .ctrls[2].Hwnd
		brdr = (GetWindowLong(hwnd, GWL.STYLE) & WS.BORDER) isnt 0 ? 1 : 0
		brdr += (GetWindowLong(hwnd, GWL.EXSTYLE) & WS_EX.CLIENTEDGE) isnt 0 ? 2 : 0
		brdr *= 2
		.scroll_child = .ctrls[2].GetChildren()[0]
		.movesplit(.r.y + .r.h - .scroll_child.Ymin - brdr - .ctrls[1].Ymin - 1)
		}
	MaximizeSecond()
		{
		.movesplit(.r.y + .ctrls[0].Ymin + 1)
		}
	moveObservers: ()
	AddMoveObserver(fn)
		{
		.moveObservers.Add(fn)
		}
	splitChanged: false
	Movesplit(n)
		{
		.splitChanged = true
		.movesplit(n)
		super.Resize(.r.x, .r.y, .r.w, .r.h)
		.callMoveObservers()
		}
	callMoveObservers()
		{
		for fn in .moveObservers
			fn()
		}
	movesplit(n)
		{
		if .Dir is "vert"
			.moveVertSplit(n)
		else // horz
			.moveHorzSplit(n)
		}
	moveVertSplit(n)
		{
		if (n < .r.y + .ctrls[0].Ymin or .ystretch0 is 0)
			.setStretch('Ystretch', 0, .Ystretch)
		else if (n > .r.y + .r.h - .ctrls[2].Ymin or .ystretch2 is 0)
			.setStretch('Ystretch', .Ystretch, 0)
		else
			{
			s = (n - .r.y - .ctrls[0].Ymin) * .Ystretch / (.r.h - .Ymin)
			// there is a rare case where this can be inf
			// if it happens ensure that it does not get saved
			// (causes issues where the split cannot be closed,
			// making the book unusable)
			if not IsInf?(s)
				.setStretch('Ystretch', Max(0, s), Max(0, .Ystretch - s))
			}
		}
	moveHorzSplit(n)
		{
		if (n < .r.x + .ctrls[0].Xmin or .xstretch0 is 0)
			.setStretch('Xstretch', 0, .Xstretch)
		else if (n > .r.x + .r.w - .ctrls[2].Xmin or .xstretch2 is 0)
			.setStretch('Xstretch', .Xstretch, 0)
		else
			{
			s = (n - .r.x - .ctrls[0].Xmin) * .Xstretch / (.r.w - .Xmin)
			// there is a rare case where this can be inf
			// if it happens ensure that it does not get saved
			// (causes issues where the split cannot be closed,
			// making the book unusable)
			if not IsInf?(s)
				.setStretch('Xstretch', Max(0, s), Max(0, .Xstretch - s))
			}
		}
	setStretch(stretch, ctrl0, ctrl2)
		{
		if ctrl2 < 0.0001 /*= handle calculation discrepancies*/
			ctrl2 = 0
		.ctrls[0][stretch] = ctrl0
		.ctrls[2][stretch] = ctrl2
		}
	GetSplit()
		{
		if (.Dir is "horz")
			return Object(.ctrls[0].Xstretch, .ctrls[2].Xstretch)
		else
			return Object(.ctrls[0].Ystretch, .ctrls[2].Ystretch)
		}
	SetSplit(n)
		{
		if (.Dir is "horz")
			{
			.ctrls[0].Xstretch = Max(0, n[0])
			.ctrls[2].Xstretch = Max(0, n[1])
			}
		else
			{
			.ctrls[0].Ystretch = Max(0, n[0])
			.ctrls[2].Ystretch = Max(0, n[1])
			}
		if not .firstResize
			.Resize(.r.x, .r.y, .r.w, .r.h)
		}
	Open()
		{ } // From HandleSplitterControl
	Close()
		{ } // From HandleSplitterControl
	Getter_CanDrag?()
		{ return true }
	CanMovesplit?(n)
		{
		if (.Dir is "vert")
			{
			if ((n < .r.y + .ctrls[0].Ymin) or (n > .r.y + .r.h - .ctrls[2].Ymin))
				return false
			}
		else
			{
			if ((n < .r.x + .ctrls[0].Xmin) or (n > .r.x + .r.w - .ctrls[2].Xmin))
				return false
			}

		return true
		}
	SaveSplit()
		{
		if .splitSaveName isnt false and .splitChanged
			UserSettings.Put(.splitSaveName $ .splitSaveNameSuffix, .GetSplit())
		}
	Destroy()
		{
		.SaveSplit()
		super.Destroy()
		}
	}
