// Copyright (C) 2017 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_one()
		{
		Assert(ScannerSeq(".Map()", {|scan| scan.Text() }) is: #('.', 'Map', '(', ')'))

		Assert(ScannerSeq("x + y",
			{|scan| scan.Type() is #WHITESPACE ? Nothing() : scan.Text() })
			is: #('x', '+', 'y'))
		}
	}