// Copyright (C) 2018 Suneido Software Corp. All rights reserved worldwide.
ScintillaAddonsControl
	{
	IDE: 			true
	addLevel: 		16
	removeLevel: 	17
	modifyLevel: 	18
	selectLevel:	19
	New(@args)
		{
		super(@.processArgs(args))
		.init()
		}

	init()
		{
		.markers = Object(
			add: .MarkerIdx(level: .addLevel)
			remove: .MarkerIdx(level: .removeLevel)
			modify: .MarkerIdx(level: .modifyLevel)
			select: .MarkerIdx(level: .selectLevel)
			)
		.indic = .IndicatorIdx(.modifyLevel)
		.IndicSetAlpha(.indic, 255/*=no transparency*/)
		.IndicSetOutlineAlpha(.indic, 255/*=no transparency*/)
		.IndicSetUnder(.indic, true)

		.indicSelect = .IndicatorIdx(.selectLevel)
		.IndicSetAlpha(.indicSelect, 255/*=no transparency*/)
		.IndicSetOutlineAlpha(.indicSelect, 255/*=no transparency*/)
		.IndicSetUnder(.indicSelect, true)
		}

	processArgs(args)
		{
		if args.showMargin
			.Addon_show_margin = true
		args.scheme = 'default'
		return args
		}

	BaseStyling()
		{
		return super.BaseStyling().Add(
			[level: .addLevel,
				marker: [SC.MARK_BACKGROUND, back: 0x00ccffcc]],
			[level: .removeLevel,
				marker: [SC.MARK_BACKGROUND, back: 0x00e6ccff]],
			[level: .modifyLevel,
				marker: [SC.MARK_BACKGROUND, back: 0x00ffe5cc]
				indicator: [INDIC.ROUNDBOX, fore: RGB(128, 187, 255)/*=blue*/]],
			[level: .selectLevel,
				marker: [SC.MARK_SHORTARROW, back: CLR.yellow],
				indicator: [INDIC.ROUNDBOX, fore: RGB(153, 201, 255)/*=blue*/]
			])
		}

	AddMarker(row, type, selected? = false)
		{
		marker = selected?
			? .markers.select
			: .markers.Member?(mem = type.Lower())
				? .markers[mem]
				: .markers.select
		.MarkerAdd(row, marker)
		}

	RemoveMarker(row)
		{
		.MarkerDelete(row, .markers.select)
		}

	AddIndic(from, length, lineStart, lineLength, selected? = false)
		{
		if selected? is true
			.SetIndicator(.indicSelect, from, length)
		else
			{
			// SetIndicatorCurrent to the type of indicator we want to remove
			.SetIndicatorCurrent(.indicSelect)
			.IndicatorClearRange(lineStart, lineLength)
			.SetIndicator(.indic, from, length)
			}
		}

	// override to disable the .SetFocus() calls
	// since it will steal the focus from the change lists in SvcControl
	GotoLine(line)
		{
		super.GotoLine(line, noFocus?: true)
		}

	SetupMargin()
		{
		.SetMarginTypeN(0, SC.MARGIN_TEXT)
		.SetMarginWidthN(0, .width = ScaleWithDpiFactor(22/*=width*/))
		}

	AddMarginText(line, text)
		{
		.MarginSetStyle(line, SC.STYLE_LINENUMBER)
		SendMessageTextIn(.Hwnd, SCI.MARGINSETTEXT, line, text)
		}

	Addon_brace_match:,
	Addon_calltips:,
	Addon_highlight_occurrences:,
	Addon_suneido_style:,
	Addon_indent_guides:,
	Addon_overwrite_lines:,
	Addon_go_to_definition:
	}
