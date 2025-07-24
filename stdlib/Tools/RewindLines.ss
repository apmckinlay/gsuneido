// Copyright (C) 2017 Suneido Software Corp. All rights reserved worldwide.
// Takes a instance of File (f), and a position (pos) to read backwards from
// numOfLines controls how many lines you read back.
//		1 = start of the current line
//		2 = start of the previous line
// 		2 > as far back as
class
	{
	CallClass(f, pos, numOfLines = 1)
		{
		f.Seek(pos)
		linesRewound = 0
		while pos > 0
			{
			c = .prevChar(f, pos)
			if c is '\n'
				linesRewound++
			if linesRewound >= numOfLines
				break
			--pos
			}
		f.Seek(pos)
		return pos
		}

	prevChar(f, pos)
		{
		f.Seek(pos - 1)
		return f.Read(1)
		}
	}