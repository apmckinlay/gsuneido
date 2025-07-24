// Copyright (C) 2010 Suneido Software Corp. All rights reserved worldwide.
function (filter, defaultControl, title? = false, except = false)
	{
	windows = Suneido.GetDefault(#Persistent, #()).GetDefault(#Windows, #())
	for window in windows
		if ((title? is true and window.Ctrl.Title is filter) or
			(title? is false and (Display(window.Ctrl.Base()) =~ filter)))
			{
			if except is window.Hwnd
				continue
			WindowActivate(window.Hwnd)
			return window.Ctrl
			}
	return PersistentWindow(defaultControl).Ctrl
	}
