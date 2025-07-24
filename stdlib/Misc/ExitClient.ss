// Copyright (C) 2022 Suneido Software Corp. All rights reserved worldwide.
function (code = 0)
	{
	if Sys.SuneidoJs?()
		SujsAdapter.CallOnRenderBackend(#Terminate, reason: '')
	else
		Exit(code)
	}