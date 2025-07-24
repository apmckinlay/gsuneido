// Copyright (C) 2001 Suneido Software Corp. All rights reserved worldwide.
// tab delimited Export
Export
	{
	Header()
		{
		.Putline(.Head.Join('\t'))
		}
	Export1(x)
		{
		line = ""
		for (field in .Fields)
			{
			// can't output newlines within the data
			s = String(x[field]).Tr("\r\n", " ")
			if (s.Has?('\t'))
				throw "ExportTab does not handle tabs within fields\r\n\r\n" $
					'\t' $ SelectPrompt(field) $ ': ' $ Display(s)
			line $= s $ '\t'
			}
		.Putline(line[.. -1])
		}
	}
