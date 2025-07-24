// Copyright (C) 2013 Suneido Software Corp. All rights reserved worldwide.
class
	{
	FromLibView(libview)
		{
		libview.Save()
		observer = .CallClass(libview.Explorer.GetTabsPaths(all?:), libview)
		libview.AlertTestResult(observer)
		}

	CallClass(paths, libview)
		{
		tests = [].Set_default([])
		for path in paths
			{
			name = LibraryTags.RemoveTagsFromName(path.AfterLast('/').Tr('?'))
			if name.Suffix?('Test')
				tests.AddUnique(name)
			else
				{
				tests.AddUnique(name $ 'Test')
				tests.AddUnique(name $ '_Test')
				}
			}

		observer = new TestObserverString(quiet:)
		if tests.Empty?()
			{
			observer.Output('No Tests Selected')
			return observer
			}

		libview.Editor.SendToAddons('On_BeforeAllTests')
		for lib in Libraries()
			{
			for name in tests.Copy()
				if .RunTest?(lib, name)
					{
					LibViewRunTest(libview.Editor, lib, name)
						{ |unused|
						CodeState.RunCurrentCode(lib, name, { .runTest(name, observer) })
						observer.ClearFailed()
						true
						}
					tests.Remove(name)
					}
			if tests.Empty?()
				break
			}
		libview.Editor.SendToAddons('On_AfterAllTests')
		return observer
		}

	RunTest?(lib, name)
		{
		if not Libraries().Has?(lib)
			return false // unused library
		if name is "" or name is "Test" and lib is "stdlib"
			return false // stdlib:Test is a helper class for tests, not a test class
		if QueryEmpty?(lib, :name, group: -1)
			return false // folder or record no longer exists
		return true
		}

	runTest(name, observer)
		{
		x = false
		resultPat = '\d tests? (SUCCEEDED|FAILED)'
		try
			x = Global(name)
		catch (e)
			if e isnt "can't find " $ name
				{
				observer.BeforeTest(name)
				observer.BeforeMethod('Global')
				observer.Error('unused', e $ '\r\n')
				observer.AfterTest('unused', 0, 0, 0)
				}
		if x isnt false and Class?(x) and x.Base?(Test)
			{
			x.RunTest(:name, :observer, quiet:)
			observer.Result = observer.Result.Replace(resultPat, name)
			}
		}
	}
