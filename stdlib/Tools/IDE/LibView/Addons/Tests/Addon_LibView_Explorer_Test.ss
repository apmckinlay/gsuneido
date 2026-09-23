// Copyright (C) 2019 Axon Development Corporation All rights reserved worldwide.
Test
	{
	Test_toggleUsed()
		{
		fn = Addon_LibView_Explorer.Addon_LibView_Explorer_toggleUsed
		Assert(fn(#test, true) is: #test)
		Assert(fn("(test)", true) is: #test)
		Assert(fn(#test, false) is: "(test)")
		Assert(fn("(test)", false) is: "(test)")
		}

	Test_RestoreTab()
		{
		mock = Mock(Addon_LibView_Explorer)
		_expectedPath = ""
		_expectedSkip = true
		mock.Explorer = FakeObject(GotoPath:
			{|path, skipFolder?|
			Assert(path is: _expectedPath)
			Assert(skipFolder? is: _expectedSkip)
			})
		mock.When.Explorer_RestoreTab([anyArgs:]).CallThrough()
		mock.When.Libs([anyArgs:]).Return(libs = [#lib0, #lib1, #lib2])

		// Should execute without error
		mock.Explorer_RestoreTab(tabOb = [path: "", group: false])

		// Lib is not in use, path should be "unused"
		_expectedPath = "(lib4)/folder/class1"
		tabOb.path = "lib4/folder/class1"
		mock.Explorer_RestoreTab(tabOb)

		// Lib is now used, path should be "used"
		libs.Add(#lib4)
		_expectedPath = "lib4/folder/class1"
		mock.Explorer_RestoreTab(tabOb)

		// Libs change, lib4 is still used, path should still be "used"
		mock.When.Libs([anyArgs:]).Return([#lib4])
		mock.Explorer_RestoreTab(tabOb)

		// Libs change, lib41 is similar but not a match, path should be "unused"
		mock.When.Libs([anyArgs:]).Return([#lib41])
		_expectedPath = "(lib4)/folder/class1"
		mock.Explorer_RestoreTab(tabOb)

		// Lib is still unused, path should still be "unused", (going to folder)
		_expectedSkip = false
		tabOb.group = true
		mock.Explorer_RestoreTab(tabOb)
		}
	}
