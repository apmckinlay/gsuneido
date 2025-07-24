// Copyright (C) 2003 Suneido Software Corp. All rights reserved worldwide.
WndProc
	{
	Name: OverviewBar
	Xmin: 7
	Ystretch: 1
	New(.priorityColor = false, .partnerCtrl = false)
		{
		.min = ScaleWithDpiFactor(2)
		.CreateWindow('SuWhiteArrow', '', WS.VISIBLE)
		.SubClass()
		.marks = Object()
		.brushes = Object()
		.brushes.bgnd = CreateSolidBrush(GetSysColor(COLOR.BTNFACE))
		}

	numRows: 0
	SetNumRows(rows)
		{
		.numRows = rows
		.Repaint()
		}

	GetNumRows()
		{
		return .numRows
		}
	SetMaxRowHeight(height, scaled? = false)
		{
		.maxRowHeight = scaled? ? height : ScaleWithDpiFactor(height)
		.Repaint()
		}

	// NOTE: code assumes that vertical scroll bar overlaps top margin
	topMargin: 0
	SetTopMargin(margin)
		{
		.topMargin = margin
		.Repaint()
		}

	AddMark(row, color = 0)
		{
		.marks[row] = color
		if not .brushes.Member?(color)
			.brushes[color] = CreateSolidBrush(color)
		.Repaint()
		}

	RemoveMark(row)
		{
		if not .marks.Member?(row)
			return

		.marks.Delete(row)
		.Repaint()
		}

	ClearMarks()
		{
		.marks = Object()
		.Repaint()
		}

	ERASEBKGND()
		{
		return 1
		}

	PAINT()
		{
		hdc = BeginPaint(.Hwnd, ps = Object())
		.paint(hdc)
		EndPaint(.Hwnd, ps)
		return 0
		}

	horzScorllBar?: false
	Scintilla_VScroll()
		{
		if .horzScorllBar? isnt HorzScrollBar?(.partnerHwnd)
			{
			.horzScorllBar? = not .horzScorllBar?
			.Repaint()
			}
		}

	top: 0
	bottom: 0
	rowHeight: 0
	rowPadding: 4
	marginPadding: 3
	paint(hdc)
		{
		GetClientRect(.Hwnd, r = Object())
		FillRect(hdc, r, GetSysColorBrush(COLOR.BTNHIGHLIGHT))
		if .numRows <= 0
			return
		height = r.bottom - r.top
		cxhscroll = GetSystemMetrics(SM.CXHSCROLL)
			// also used for height of vertical scroll arrows
			// on the assumption they are "square"

		.horzScorllBar? = HorzScrollBar?(.partnerHwnd)
		hsAllow = .horzScorllBar? ? cxhscroll : 0

		// assume not scrolling
		.top = .topMargin
		.rowHeight = .maxRowHeight
		.bottom = height - (.top + .numRows * .rowHeight + .rowPadding)

		if .bottom < hsAllow
			{
			// wrong - must be vertical scroll bar
			// need to adjust margins
			// so the bar matches the scrolling thumb range
			if .topMargin is 0
				.top = cxhscroll
			.bottom = cxhscroll + hsAllow // bottom scroll arrow + horz scroll bar
			.rowHeight = (height - (.top + .bottom) - .rowPadding) / .numRows
			}

		DoWithHdcObjects(hdc, [GetStockObject(SO.NULL_PEN), .brushes.bgnd])
			{
			Rectangle(hdc, // top margin
				r.left, r.top, r.right + 1, r.top + .top + .marginPadding)
			Rectangle(hdc, // bottom margin
				r.left, r.bottom - .bottom - 2, r.right + 1, r.bottom + 1)
			.paintMarks(hdc, r)
			}
		.bottom = r.bottom - .bottom
		}

	getter_partnerHwnd()
		{
		if .partnerCtrl isnt false
			return .partnerHwnd = .partnerCtrl.Hwnd

		siblings = .Parent.GetChildren()
		if false is (i = siblings.Find(this)) or
			i is 0 or not siblings[i - 1].Member?('Hwnd')
			return 0
		return .partnerHwnd = siblings[i - 1].Hwnd // once only
		}

	paintMarks(hdc, r)
		{
		respectPriority? = .rowHeight < 1 and .priorityColor isnt false
		overwriteRate = (.min / .rowHeight).Round(0)
		overwrite = 0
		for row in .marks.Members().Sort!()
			{
			color = overwrite-- > 0 ? .priorityColor : .marks[row]
			if respectPriority? and .marks[row] is .priorityColor
				overwrite = overwriteRate
			SelectObject(hdc, .brushes[color]) // restored via calling code
			top = .top + 2 + .rowHeight * row
			Rectangle(hdc, r.left, top, r.right + 1, top + Max(.rowHeight, .min) + 1)
			}
		}

	LBUTTONUP(lParam)
		{
		y = HIWORD(lParam)
		if .rowHeight is 0 or y < .top or y > .bottom
			return 0
		row = ((y - .top) / .rowHeight).Floor()
		.Send('Overview_Click', row)
		return 0
		}

	Destroy()
		{
		.brushes.Each(DeleteObject)
		super.Destroy()
		}
	}
