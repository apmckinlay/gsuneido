// Copyright (C) 2001 Suneido Software Corp. All rights reserved worldwide.
// CSV Export
Export
	{
	Ext: "csv"
	Header()
		{
		.Putline(.Head.JoinCSV())
		}
	Export1(record)
		{
		// can't output newlines within the data
		.Putline(record.JoinCSV(.Fields).Tr("\r\n", " "))
		}
	}
