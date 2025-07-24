// Copyright (C) 2000 Suneido Software Corp. All rights reserved worldwide.
RowFormat
	{
	Print(@args/*unused*/)
		{
		}
	GetSize(data/*unused*/ = #())
		{ return #(w: 0, h: 0, d: 0); }
	ExportCSV(data/*unused*/ = '')
		{
		return ''
		}
	}
