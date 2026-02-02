// Copyright (C) 2025 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_customizeScreen?()
		{
		m = CustomizeDialog.CustomizeDialog_customizeScreen?

		Assert(m(true, false, false))
		Assert(m(false, true, false))
		Assert(m(false, false, Object(custom_tabs: #(Header))))

		Assert(m(false, false, Object(custom_tabs: #())) is: false)
		Assert(m(false, false, Object()) is: false)
		Assert(m(false, false, false) is: false)
		}
	}