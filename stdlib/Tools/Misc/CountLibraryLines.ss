// Copyright (C) 2013 Suneido Software Corp. All rights reserved worldwide.
function ()
	{
	total_records = total_lines = 0
	recordsColWidth = 6
	linesColWidth = 7
	for lib in LibraryTables()
		{
		records = lines = 0
		QueryApply(lib, group: -1)
			{|x|
			++records
			lines += x.lib_current_text.LineCount()
			}
		Print(records.Pad(recordsColWidth, ' '), lines.Pad(linesColWidth, ' '), lib)
		total_records += records
		total_lines += lines
		}
	Print("====== =======")
	Print(total_records.Pad(recordsColWidth, ' '), total_lines.Pad(linesColWidth, ' '))
	}
