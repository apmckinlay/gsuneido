// Copyright (C) 2026 Suneido Software Corp. All rights reserved worldwide.
_ScintillaDiffControl
	{
	ComponentName: #ScintillaDiff

	init()
		{
		.markers = Object(
			add: .MarkerIdx(level: .addLevel)
			remove: .MarkerIdx(level: .removeLevel)
			modify: .MarkerIdx(level: .modifyLevel)
			select: .MarkerIdx(level: .selectLevel)
			)
		.indic = .IndicatorIdx(.modifyLevel)
		.indicSelect = .IndicatorIdx(.selectLevel)
		}

	AddMarginText(line, text)
		{
		SuServerPrint('AddMarginText', line, text)
//		.MarginSetStyle(line, SC.STYLE_LINENUMBER)
//		SendMessageTextIn(.Hwnd, SCI.MARGINSETTEXT, line, text)
		}
	}