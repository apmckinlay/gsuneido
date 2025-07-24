// Copyright (C) 2017 Suneido Software Corp. All rights reserved worldwide.
function()
	{
	if not Suneido.Member?("QcEnabled")
		{
		Suneido.QcEnabled = UserSettings.Get("QcEnabled", true)
		UserSettings.AddObserver("QcEnabled", { |value| Suneido.QcEnabled = value})
		}
	return Suneido.QcEnabled
	}