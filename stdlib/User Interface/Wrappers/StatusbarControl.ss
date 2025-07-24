// Copyright (C) 2000 Suneido Software Corp. All rights reserved worldwide.
WndProc
	{
	Name:		"Status"
	Xstretch:	1
	New()
		{
		.CreateWindow(STATUSCLASSNAME, NULL, WS.VISIBLE)
		.SetFont()
		r = GetWindowRect(.Hwnd)
		.Ymin = r.bottom - r.top
		}

	Resize(x, y, w, h)
		{
		if not .Member?('Hwnd')
			return
		MoveWindow(.Hwnd, x, y , w, h, true)
		}

	Set(text, panel = 0)
		// pre:		text is a string AND
		//			0 <= panel < .GetNumPanels OR panel is false
		// post:	if panel is an integer, the statusbar part indexed by panel's
		//			text string is text
		//			otherwise if panel is false, the simple panel's is set to text
		{
		Assert((panel is false) or (0 <= panel and panel < .GetNumPanels()))
		Assert(String?(text))
		text = TranslateLanguage(text)
		if panel is false		// set simple panel
			SendMessageTextIn(.Hwnd, SBM.SETTEXT,
				255 /*= SB_SIMPLEID */ | .GetPanelStyle(panel), text)
		else					// set indexed panel
			SendMessageTextIn(.Hwnd, SBM.SETTEXT, panel | .GetPanelStyle(panel), text)
		}

	Get(panel = 0)
		// pre:		0 <= panel < .GetNumPanels
		// post:	returns text associated with the statusbar part indexed by panel
		{
		Assert(0 <= panel and panel < .GetNumPanels())
		len = LOWORD(.SendMessage(SBM.GETTEXTLENGTH, panel, 0))
		return SendMessageTextOut(.Hwnd, SBM.GETTEXT, panel, len + 1).text
		}

	SetPanelStyle(panel, style)
		// pre:		0 <= panel < .GetNumPanels AND
		// 			style is valid style constant (SBT)
		// post:	the statusbar part indexed by panel has style as its style
		{
		Assert(panel < .GetNumPanels())
		text = .Get(panel)
		SendMessageTextIn(.Hwnd, SBM.SETTEXT, panel | style, text)
		}

	GetPanelStyle(panel)
		// pre:		0 <= panel < .GetNumPanels
		// post:	the style of the statusbar part indexed by panel is returned
		{
		Assert(panel < .GetNumPanels())
		return HIWORD(.SendMessage(SBM.GETTEXTLENGTH, panel, 0))
		}

	SetBkColor(color = false)
		// pre:		color is an integer value (COLORREF)
		// post:	this' background color is color
		{
		if color isnt false
			.DisableTheme() // themes don't support SETBKCOLOR
		.SendMessage(SB.SETBKCOLOR, 0, color is false ? CLR.DEFAULT : color)
		}

	AddPanel(size = 100, at = false, style = false)
		// pre:		size is a positive integer AND
		//			at is false OR 0 <= at <= .GetNumPanels AND
		//			style is false OR style is a valid style constant (SBT)
		{
		Assert((at is false) or (0 <= at and at <= .GetNumPanels()))
		newPanels = .GetPanels()
		if at isnt false
			newPanels = .addNewPanels(newPanels, at, size)
		else
			{
			right = newPanels.Empty?() ? size : size + newPanels.Last()
			newPanels.Add(right)
			}
		.SetPanels(newPanels)
		if style isnt false
			.SetPanelStyle(at isnt false ? at : .GetNumPanels() - 1, style)
		}

	addNewPanels(newPanels, at, size)
		{
		if not newPanels.Empty?()
			{
			right = at is 0
				? size
				: size + newPanels[at - 1]
			newPanels.Add(right, :at)
			numPanels = .GetNumPanels()
			for (i = at + 1; i <= numPanels; i++)
				newPanels[i] += size
			}
		else
			newPanels = Object(size)
		return newPanels
		}

	RemPanel(at = false)
		// pre:		at is false OR
		//			0 <= at < .GetNumPanels
		// post:	if at is false, (pre)last panel has been removed
		//			else (pre).GetPanels()[at] has been removed
		{
		Assert((at is false) or (0 <= at and at < .GetNumPanels()))
		if at is false
			at = .GetNumPanels() - 1
		styles = Object()
		panels = .GetPanels()
		.shiftPanels(panels, at, styles)
		panels.Delete(at)
		.SetPanels(panels)
		for style in styles.Members()
			.SetPanelStyle(style, styles[style])
		}

	shiftPanels(panels, at, styles)
		{
		size = at is 0 ? panels[at] : panels[at] - panels[at - 1]
		for panel in panels.Members()
			{
			if panel isnt at
				styles.Add(.GetPanelStyle(panel))
			if panel > at
				panels[panel] -= size
			}
		}

	GetNumPanels()
		// post:	the number of panels contained by this is returned
		{ return .SendMessage(SB.GETPARTS, 0, 0) }

	GetPanels()
		// post:	an object containing a list of integers representing the right-
		//			hand sides of every panel is returned
		{
		num = SendMessageSBPART(.Hwnd, SB.GETPARTS, .GetNumPanels(), panels = Object())
		result = Object()
		while --num >= 0
			result.Add(panels.parts[num], at: 0)
		return result is Object(-1) ? Object() : result
		}

	SetPanels(panels = #(-1))
		// pre:		panels is a list of integers representing the right-hand side of
		//			every panel to be added
		// post:	this contains panels.Size() statusbar parts with right-hand-sides
		//			corresponding to their index in panels
		{
		count = panels.Size()
		Assert(count < 256 /*= max number of panels*/)
		return SendMessageSBPART(.Hwnd, SB.SETPARTS, count, Object(parts: panels))
		}

	GetSimple()
		// post:	returns true iff this' statusbar is in simple mode
		{ return .SendMessage(SB.ISSIMPLE, 0, 0) isnt 0 }

	SetSimple(simple = true)
		// pre:		simple is a boolean value
		// post:	this is in simple mode iff simple is true
		{ .SendMessage(SB.SIMPLE, simple, 0) }

	GetPanelRect(panel)
		// pre:		0 <= panel < .GetNumPanels
		// post:	returns a Rect object containing the coordinates of the rectangle
		//			of the statusbar part indexed by panel
		{
		Assert(0 <= panel and panel < .GetNumPanels())
		SendMessageRect(.Hwnd, SB.GETRECT, panel, rc = Object())
		return Rect(rc.left, rc.top, rc.right - rc.left, rc.bottom - rc.top)
		}

	GetReadOnly()			// read-only not applicable to statusbar
		{ return true }
	}
