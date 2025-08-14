// Copyright (C) 2013 Suneido Software Corp. All rights reserved worldwide.
function (libview, lib, name)
	{
	if not RunAssociatedTests.RunTest?(lib, name)
		return
	name = LibraryTags.RemoveTagFromName(name)
	libview.Save()
	libview.Editor.SendToAddons('On_BeforeAllTests')
	CodeState.RunCurrentCode(lib, name)
		{
		x = Global(name)
		if Class?(x) and x.Base?(Test)
			if LibViewRunTest(libview.Editor, lib, name, { it.Debug() })
				{
				libview.Editor.SendToAddons('On_AfterAllTests')
				libview.AlertInfo("Run Test", name $ " SUCCEEDED")
				}
		}
	}

