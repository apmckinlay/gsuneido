// Copyright (C) 2004 Suneido Software Corp. All rights reserved worldwide.
WndProc
	{
	Name: ReporterColHeads
	Xmin: 700
	New()
		{
		.CreateWindow("SuWhiteArrow", "", WS.VISIBLE)
		.SubClass()
		.report = ReportInstance(fontScale: 20)
		}
	colheads: false
	SetColumns(cols, widths)
		{
		.cols = cols
		.widths = widths
		.setColumns()
		.Repaint()
		}
	setColumns()
		{
		WithDC(.Hwnd)
			{ |dc|
			.setScale(dc)
			.report.SetDC(dc)
			_report = .report
			headings = .cols.Map({ Object(Heading: it) })
			scale = 11.InchesInTwips() / .Xmin
			widths = .widths.Map({ it * scale - HskipFormat.Size.w })
			.colheads = ColHeadsFormat(headings, widths,
				#(name: Arial, size: 10, weight: 400))
			.height = .colheads.GetSize().h
			.Ymin = .height / scale
			}
		.Window.Refresh()
		}
	PAINT()
		{
		if .colheads is false
			return 0
		dc = BeginPaint(.Hwnd, ps = Object())
		.setScale(dc)
		GetClientRect(.Hwnd, r = Object())
		.report.SetDC(dc)
		_report = .report
		r.left += 100
		// use 100000 for width to ensure HorzFormat does not scale the header items
		.colheads.Print(r.left, r.top, 100000, r.bottom - r.top)
		EndPaint(.Hwnd, ps)
		return 0
		}
	w: 200
	setScale(dc)
		{
		SetMapMode(dc, MM.ISOTROPIC)
		SetWindowExt(dc, 11.InchesInTwips(), 11.InchesInTwips())
		SetViewportExt(dc, .Xmin, .Xmin)
		}
	Destroy()
		{
		.report.Destroy()
		super.Destroy()
		}
	}