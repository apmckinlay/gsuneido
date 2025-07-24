// Copyright (C) 2006 Suneido Software Corp. All rights reserved worldwide.
// e.g. TrackBarControl(range: [0, 100])
// see http://www.codeguru.com/Cpp/controls/statusbar/article.php/c2969
WndProc
	{
	Name: "TrackBar"

	New(vert = false, noticks = false, tickmarks = 'BOTTOM', enableselrange = false,
		fixedlength = true, nothumb = false, slidetip = true, style = 0,
		range = false, start = false, ticfreq = false, tip = "")
		{
		if vert
			{
			style |= TBS.VERT
			if .Xmin is 0
				.Xmin = ScaleWithDpiFactor(30)
			if .Ymin is 0
				.Ymin = ScaleWithDpiFactor(100)
			}
		else
			{
			if .Xmin is 0
				.Xmin = ScaleWithDpiFactor(100)
			if .Ymin is 0
				.Ymin = ScaleWithDpiFactor(30)
			}
		style |= noticks ? TBS.NOTICKS : TBS.AUTOTICKS | TBS[tickmarks]
		if enableselrange
			style |= TBS.ENABLESELRANGE
		if fixedlength
			style |= TBS.FIXEDLENGTH
		if nothumb
			style |= TBS.NOTHUMB
		if slidetip
			style |= TBS.TOOLTIPS
		.CreateWindow('msctls_trackbar32', '', WS.VISIBLE | style)
		if range isnt false
			.SetRange(range[0], range[1])
		if ticfreq isnt false
			.SetTicFreq(ticfreq[0], ticfreq[1])
		if start isnt false
			.Set(start)
		.SetTip(tip)
		}

	HSCROLL(wParam)
		{
		.isendtrack?(wParam)
		return 0
		}
	VSCROLL(wParam)
		{
		.isendtrack?(wParam)
		return 0
		}
	isendtrack?(wParam)
		{
		// add to TB: TB.ENDTRACK: 0x0008
		if LOWORD(wParam) is TB.ENDTRACK
			.Send("NewValue", .Get())
		//if (dirty?)
		}
	SetTip(tip)
		{
		.ToolTip(TranslateLanguage(tip))
		}

	SetRange(min, max, redraw = true)
		{
		.SendMessage(TBM.SETRANGE, redraw, MAKELONG(min, max))
		}
	Set(position)
		{
		.SendMessage(TBM.SETPOS, 1, position)
		}
	Get()
		{
		return .SendMessage(TBM.GETPOS)
		}

	ClearTics(redraw = true)
		{
		.SendMessage(TBM.CLEARTICS, redraw, 0)
		}
	SetTic(position)
		{
		return .SendMessage(TBM.SETTIC, 0, position) is 1
		}
	SetTicFreq(freq, position)
		{
		.SendMessage(TBM.SETTICFREQ, freq, position)
		}

	ClearSel(redraw = true)
		{
		.SendMessage(TBM.CLEARSEL, redraw, 0)
		}
	SetSel(start, end, redraw = true)
		{
		.SendMessage(TBM.SETSEL, redraw, MAKELONG(start, end))
		}
	GetSelStart()
		{
		return .SendMessage(TBM.GETSELSTART)
		}
	GetSelEnd()
		{
		return .SendMessage(TBM.GETSELEND)
		}
	}
