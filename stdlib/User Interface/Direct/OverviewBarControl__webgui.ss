// Copyright (C) 2026 Suneido Software Corp. All rights reserved worldwide.
Control
	{
	Name: #OverviewBar
	ComponentName: #OverviewBar
	New(priorityColor = false, .partnerCtrl = false)
		{
		.ComponentArgs = Object(priorityColor)
		if .partnerCtrl isnt false
			.ComponentArgs.Add(.partnerCtrl.Hwnd)
		}

	numRows: 0
	SetNumRows(rows)
		{
		.numRows = rows
		.Act(#SetNumRows, rows)
		}

	GetNumRows()
		{
		return .numRows
		}

	SetMaxRowHeight(that, method, scaled?/*unused*/ = false)
		{
		.Act(#SetMaxRowHeight, that.Hwnd, method)
		}

	SetTopMargin(unused)
		{
		// FIX me
		}

	AddMark(row, color = 0)
		{
		.Act(#AddMark, row, color)
		}

	RemoveMark(row)
		{
		.Act(#RemoveMark, row)
		}

	ClearMarks()
		{
		.Act(#ClearMarks)
		}

	Default(@args)
		{
		SuServerPrint(args[0])
		}
	}