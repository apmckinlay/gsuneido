// Copyright (C) 2000 Suneido Software Corp. All rights reserved worldwide.
// print a library record
// generates a CodeFormat for each line of code
Generator
	{
	New(.name, text)
		{
		.input = text.Lines().Iter()
		}
	Header()
		{
		return Object('Vert',
			Object("Text", .name, font: #(name: "Arial", size: 14, weight: 'bold')),
			'Vskip')
		}
	Next()
		{
		if .input is line = .input.Next()
			return false
		return _report.Construct(Object('Vert'
			Object("Code", line, w: _report.GetWidth())
			#(Vskip .02)))
		}
	}
