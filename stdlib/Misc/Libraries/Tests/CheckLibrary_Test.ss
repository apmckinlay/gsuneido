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

	Test_suppressions()
		{
		.MakeLibraryRecord(
			[name: 'Test_CheckLibrary_HasSuppressed',
				text: "function () { First_NonExistingFn() }"],
			[name: 'Test_CheckLibrary_NotSuppressed',
				text: "function () { Second_NonExistingFn() }"],
			[name:  'Test_lib_CheckLibrarySuppressions',
				text: `#(Test_CheckLibrary_HasSuppressed)`])
		ck = CheckLibrary("Test_lib")
		Assert(ck has: 'Test_lib:Test_CheckLibrary_NotSuppressed')
		Assert(ck hasnt: 'Test_CheckLibrary_HasSuppressed')
		}
	}
