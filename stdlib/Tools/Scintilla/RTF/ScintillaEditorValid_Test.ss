// Copyright (C) 2025 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_one()
		{
		Assert(ScintillaEditorValid('', false))
		Assert(ScintillaEditorValid('test', false))

		Assert(ScintillaEditorValid('1234567890', 5) is: false)
		Assert(ScintillaEditorValid('1234567890', 10))

		Assert(ScintillaEditorValid('1234567890x', 10) is: false)
		}
	}