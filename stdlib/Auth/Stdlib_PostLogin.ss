// Copyright (C) 2024 Suneido Software Corp. All rights reserved worldwide.
function ()
	{
	if Sys.Win32?()
		Thread(WebView2.CleanUp)
	}