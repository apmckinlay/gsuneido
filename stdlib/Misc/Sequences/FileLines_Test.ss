// Copyright (C) 2024 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_one()
		{
		file = .MakeFile(s = 'one\ntwo\nthree')
		t = FileLines(file) { it.Join('\n') }
		Assert(t is: s)
		}
	}