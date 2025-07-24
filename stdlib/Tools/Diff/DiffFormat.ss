// Copyright (C) 2000 Suneido Software Corp. All rights reserved worldwide.
Generator
	{
	New(code1, code2, name = "")
		{
		.diffs = Diff.SideBySide(code1.Lines(), code2.Lines())
		.name = name
		.i = -1
		}
	Header()
		{
		return Object('Vert',
			Object("Text", .name, font: #(name: "Arial", size: 14)),
			'Vskip')
		}
	Next()
		{
		++.i
		if (.i >= .diffs.Size())
			return false
		width = (_report.GetWidth() - 720 /*= half an inch in Twips*/) / 2
		return _report.Construct(Object('Horz'
			Object('Code', .diffs[.i][0], w: width, fontsize: 8)
			'Hskip' Object('Text', .diffs[.i][1], width: 1) 'Hskip'
			Object('Code', .diffs[.i][2], w: width, fontsize: 8)))
		}
	}
