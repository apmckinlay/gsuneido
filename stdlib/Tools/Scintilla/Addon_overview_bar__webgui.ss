// Copyright (C) 2026 Suneido Software Corp. All rights reserved worldwide.
_Addon_overview_bar
	{
	Overview_Reset()
		{
		if .ovbars.Empty?()
			return

		.adjust()
		ovbarHwnds = Object()
		for type in .ovbars.Members()
			ovbarHwnds[type] = .ovbars[type].Hwnd
		.Act(#UpdateOverview, ovbarHwnds, .markersInfo)
		}

	getter_markersInfo()
		{
		markersInfo = Object().Set_default(Object())
		for type in .GetMarkerTypes()
			.ForEachMarkerByLevel(type)
				{ |idx|
				markersInfo[type].Add(idx)
				}
		return .markersInfo = markersInfo
		}

	Overview_Click(row)
		{
		.SetFirstVisibleLine(row, centerInScreen?:)
		.GotoLine(row)
		}
	}