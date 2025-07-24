// Copyright (C) 2000 Suneido Software Corp. All rights reserved worldwide.
// see also: LoadText
function (table)
	{
	File(table $ ".txt", 'w')
		{ |f|
		f.Writeline(Schema(table).AfterFirst("\n").Tr("\r\n\t", ""))
		QueryApply(table)
			{ |x|
			f.Writeline(Display(x))
			}
		}
	}
