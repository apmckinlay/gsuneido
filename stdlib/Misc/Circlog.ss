// Copyright (C) 2015 Suneido Software Corp. All rights reserved worldwide.
class
	{
	New(.size = 100)
		{
		.log = Object()
		.i = 0
		}
	Append(s)
		{
		.log[.i] = s
		.i = (.i + 1) % .size
		}
	ToString()
		{
		return Opt(.log[.i ..].Join('\n'), '\n') $
				Opt(.log[.. .i].Join('\n'), '\n')
		}
	}