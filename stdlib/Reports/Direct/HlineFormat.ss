// Copyright (C) 2000 Suneido Software Corp. All rights reserved worldwide.
Format
	{
	Xstretch: 1
	New(.width = 0, .thick = 10, .before = 120, .after = 120, .descent = 0,
		.noline = false, export = false)
		{
		.Export = export isnt false
		.double? = export is 'double'
		}
	GetSize(data /*unused*/ = false)
		{
		return Object(w: .width, h: .before + .thick + .after, d: .descent)
		}
	Print(x, y, w, h /*unused*/, data /*unused*/ = false)
		{
		if .noline is true
			return

		y += .before
		_report.AddLine(x, y, x + w, y, .thick)
		}

	Export: false
	double?: false
	ExportCSV(data /*unused*/ = false)
		{
		return .CSVExportString(.double? is false
			// use more sigle dash to match the length of double line on spreadsheet
			? '----------------------------------------------'
			// use single quote so spreadsheet does not treat it as formula
			: "'=======================")
		}
	}