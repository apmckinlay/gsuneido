// Copyright (C) 2009 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_main()
		{
		for c in #('yes', true, 'true', 'y', 'Yes', 'YES')
			Assert(Field_boolean.Encode(c))
		for c in #('no', 'n', 'N', 'false', false)
			Assert(Field_boolean.Encode(c) is: false)
		for c in #(true, false, 'xxx', 123)
			Assert(Field_boolean.Encode(c) is: c)
		}
	}