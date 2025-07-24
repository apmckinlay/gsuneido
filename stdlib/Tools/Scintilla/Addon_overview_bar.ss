// Copyright (C) 2020 Suneido Software Corp. All rights reserved worldwide.
CodeViewAddon
	{
	Name: 	OverviewBars
	Inject: editor

	ovbars: #()
	InjectControls(container)
		{
		bars = Object()
		types = .GetMarkerTypes().Sort!()
		for type in types
			bars.Add([#OverviewBar, partnerCtrl: .Parent, name: #ovbar_ $ type])
		container.AppendAll(bars)

		.ovbars = Object()
		for type in types
			.ovbars[type] = container.FindControl(#ovbar_ $ type)
		}

	AddonReady?()
		{ return .ovbars.NotEmpty?() }

	Addon_RedirMethods()
		{ return #(Overview_Click) }

	Init()
		{
		.Overview_Reset()
		}

	adjust()
		{
		for ovbar in .ovbars
			{
			ovbar.SetNumRows(.GetLineCount())
			ovbar.SetMaxRowHeight(.TextHeight(), scaled?:)
			}
		}

	Scintilla_VScroll()
		{
		if .ovbars.NotEmpty?()
			.ovbars.Each(#Scintilla_VScroll)
		}

	Overview_Click(row)
		{
		lines = .LinesOnScreen()
		.SetFirstVisibleLine(row - (lines / 2))
		.GotoLine(row)
		}

	Overview_Reset()
		{
		if .ovbars.NotEmpty?()
			.ovbars = .reset(.ovbars.Copy())
		}

	reset(ovbars)
		{
		.adjust()
		ovbars.Each(#ClearMarks)
		for i in .GetMarkerLines()
			{
			markers = .MarkerGet(i)
			for type in .GetMarkerTypes()
				.ForEachMarkerByLevel(type)
					{ |j|
					stop = false
					if 0 isnt (markers & (1 << j))
						{
						ovbars[type].AddMark(i, .GetMarkerColor(j))
						stop = true
						}
					stop
					}
			}

		if .curLine isnt false
			ovbars.Each({ it.AddMark(.curLine, CLR.GRAY) })

		return ovbars
		}

	Scintilla_MarkersChanged()
		{ .Overview_Reset() }

	curLine: false
	Scintilla_Selection()
		{
		line = .LineFromPosition(.GetCurrentPos())
		if line is .curLine
			return
		.curLine = line
		.Overview_Reset()
		}

	Scintilla_Zoom()
		{
		// need delay or you get previous text height
		.Defer(.adjust, uniqueID: #scintilla_zoom)
		}
	}
