// Copyright (C) 2019 Axon Development Corporation All rights reserved worldwide.
Test
	{
	Test_toggleUsed()
		{
		fn = Addon_LibView_Explorer.Addon_LibView_Explorer_toggleUsed
		Assert(fn('test', 	true)  is: 'test')
		Assert(fn('(test)', true)  is: 'test')
		Assert(fn('test', 	false) is: '(test)')
		Assert(fn('(test)', false) is: '(test)')
		}

	Test_RestoreTab()
		{
		mock = Mock(Addon_LibView_Explorer)
		_expectedPath = ''
		mock.Explorer = FakeObject(GotoPath: { | path | Assert(path is: _expectedPath)})
		mock.When.Explorer_RestoreTab([anyArgs:]).CallThrough()
		mock.When.Libs([anyArgs:]).Return(libs = ['lib0', 'lib1', 'lib2'])

		// Should execute without error
		mock.Explorer_RestoreTab('')

		// Lib is not in use, path should be "unused"
		_expectedPath = '(lib4)/folder/class1'
		mock.Explorer_RestoreTab(path = 'lib4/folder/class1')

		// Lib is now used, path should be "used"
		libs.Add('lib4')
		_expectedPath = 'lib4/folder/class1'
		mock.Explorer_RestoreTab(path)

		// Libs change, lib4 is still used, path should still be "used"
		mock.When.Libs([anyArgs:]).Return(['lib4'])
		mock.Explorer_RestoreTab(path)

		// Libs change, lib41 is similar but not a match, path should be "unused"
		mock.When.Libs([anyArgs:]).Return(['lib41'])
		_expectedPath = '(lib4)/folder/class1'
		mock.Explorer_RestoreTab(path)

		// Lib is still unused, path should still be "unused"
		mock.Explorer_RestoreTab('(lib4)/folder/class1')
		}
	}