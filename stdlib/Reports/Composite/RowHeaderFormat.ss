// Copyright (C) 2001 Suneido Software Corp. All rights reserved worldwide.
RowFormat
	{
	colheads: false
	New(@args)
		{
		super(@args)
		.colheads = _report.Construct(.Header())
		}
	Print(@args)
		{ .colheads.Print(@args) }
	GetSize(@args)
		{ return .colheads.GetSize(@args) }
	ExportCSV(data = '')
		{ .colheads.ExportCSV(data) }
	}
