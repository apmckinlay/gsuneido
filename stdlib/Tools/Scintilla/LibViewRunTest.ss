// Copyright (C) 2021 Suneido Software Corp. All rights reserved worldwide.
function (editor, lib, name, block)
	{
	editor.SendToAddons('On_BeforeTest', :lib, :name)
	result = TestRunner.RunTests({ block(Global(name)) })
	editor.SendToAddons('On_AfterTest', :lib, :name)
	return result
	}
