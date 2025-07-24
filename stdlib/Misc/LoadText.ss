// Copyright (C) 2000 Suneido Software Corp. All rights reserved worldwide.
// see also: DumpText
function (table)
	{
	File(table $ ".txt", 'r')
		{ |f|
		if (false is (line = f.Readline()))
			return false
		try Database("destroy " $ table)
		Database("create " $ table $ " " $ line)
		Transaction(update:)
			{ |t|
			q = t.Query(table)
			while (false isnt line = f.Readline())
				{
				while (not line.Suffix?(']') and
					false isnt morelines = f.Readline())
					line $= '\n' $ morelines
				q.Output(line.SafeEval())
				}
			q.Close()
			}
		}
	return true
	}
