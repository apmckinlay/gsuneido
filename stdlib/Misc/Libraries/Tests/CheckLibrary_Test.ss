// Copyright (C) 2013 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_old_update?()
		{
		f = CheckLibrary.CheckLibrary_old_update?
		d = CheckLibrary.CheckLibrary_old_days
		Assert(f(#(name: 'xxx')) is: false)
		Assert(f(#(name: 'Update_20010203')))
		Assert(f([name: 'Update_' $ Date().Format('yyyyMMdd')]) is: false)
		tomorrow = Date().Minus(days: d + 1)
		Assert(f([name: 'Update_' $ tomorrow.Format('yyyyMMdd')]))
		yesterday = Date().Minus(days: d - 1)
		Assert(f([name: 'Update_' $ yesterday.Format('yyyyMMdd')]) is: false)
		}

	Test_builtdate_comment()
		{
		f = CheckLibrary.BuiltDate_skip?
		Assert(f("stuff") is: false)
		Assert(f("// BuiltDate > 20010203") is: false)
		Assert(f("// BuiltDate > 20990203"))
		}
	}
