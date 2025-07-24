// Copyright (C) 2018 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	editor: class
		{
		s: "\r\n\n  \t \r\n" $ "\ttest1\r\n" $ "\r\n" $ "test2 test3\r\n \t test4\r\n" $
			" \r\n   "
		Get(_s = false)
			{
			return s isnt false ? s : .s
			}
		GetSelText(_sel = "")
			{
			return sel
			}
		GetCurrentPos(_pos)
			{
			return pos
			}
		SetSelect(start, len, _expected = false)
			{
			if expected isnt false
				Assert(.s[start::len].Trim() is: expected)
			}
		}
	Test_one()
		{
		.testGetSelText("test1", "test1")
		.testGetSelText("test1\r\n", "test1")

		.testGet("", s: "  \t\r\n")

		// clicking "\r\n\n  \t \r\n"
		for pos in ..9
			.testGet("", :pos)
		// clicking "\ttest1\r\n"
		for pos in 9..17
			.testGet("test1", :pos)
		// clicking "\r\n"
		for pos in 17..19
			.testGet("", :pos)
		// clicking "test2 test3\r\n \t test4\r\n"
		for pos in 19..42
		.testGet("test2 test3\r\n \t test4", :pos)
		// clicking " \r\n   "
		for pos in 42..48
			.testGet("", :pos)
		}

	testGetSelText(sel, expected)
		{
		_sel = sel
		Assert(GetRunText(.editor).Trim() is: expected)
		}

	testGet(expected, pos = false, s = false)
		{
		if s isnt false
			_s = s
		if pos isnt false
			_pos = pos
		_expected = expected
		Assert(GetRunText(.editor).Trim() is: expected)
		}
	}
